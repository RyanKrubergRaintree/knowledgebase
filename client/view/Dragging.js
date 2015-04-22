'use strict';

export var ItemType = "text/fedwiki";
export var IdType = "text/fedwiki+id";
export var SizeType = "text/fedwiki+height";

//TODO: placeholder keeps lingering around!!!
export class DragContext {
	constructor(proxy){
		this.proxy = proxy;
		this.item = null;
		this.node = null;
		this.hidden = false;
		this.startVersion = proxy.page.version;
	}

	get page(){
		if(this.proxy){
			return this.proxy.page;
		}
		return undefined;
	}

	start(ev, item, node, created){
		var page = this.page;
		this.item = item;
		this.node = node;

		if(!page || page.readonly || created){
			ev.dataTransfer.effectAllowed = 'copy';
		} else {
			ev.dataTransfer.effectAllowed = 'all';
		}

		var off = mouseOffset(ev);
		if(node != null){
			ev.dataTransfer.setDragImage(node, off.x, off.y);
		}

		var data = {
			item: item,
			created: created
		};
		if(page){
			data.url = page.url;
			data.title = page.title;
		}

		ev.dataTransfer.setData(ItemType, JSON.stringify(data));

		if(node){
			ev.dataTransfer.setData(SizeType, node.clientHeight);
		}
		ev.dataTransfer.setData(IdType, item.id);
	}

	drag(ev){
		if(ev.dataTransfer.dropEffect == "move"){
			this.hide();
		} else {
			this.show();
		}
		ev.preventDefault();
	}

	end(ev){
		this.show();
		if((this.startVersion == this.proxy.page.version) &&
			(ev.dataTransfer.dropEffect == "move")){
			this.proxy.modify({
				type: "remove",
				id: this.item.id
			});
		}

		this.item = null;
		this.node = null;

		ev.preventDefault();
	}

	hide(){
		this.hidden = true;
		if(this.node === null){
			return;
		}

		if(this.hidden == this.node.classList.contains("item-dragging")){
			return;
		}

		var self = this;
		window.setTimeout(function(){
			this.node.classList.toggle('item-dragging', self.hidden);
		}, 0);
	}

	show(){
		this.hidden = false;
		if(this.node === null){
			return;
		}
		this.node.classList.toggle('item-dragging', this.hidden);
	}
}

export class DropArea {
	constructor(proxy, onGetContainer, onEditItem){
		this.proxy = proxy;
		this.onGetContainer = onGetContainer;
		this.onEditItem = onEditItem;
		this.placeholder = new Placeholder();
	}

	get page(){ return this.proxy.page; }

	dropEffectFor(ev){
		if(ev.dataTransfer.effectAllowed === 'copy') {
			return 'copy';
		}
		if(ev.shiftKey){
			return 'copy';
		}
		return 'move';
	}

	enter(ev){

	}

	leave(ev){
		var container = this.onGetContainer();
		if(!container.contains(ev.target)){
			this.placeholder.forget();
		}
		if(ev.target == container){
			this.placeholder.forget();
		}
	}

	drop(ev){
		this._drop(ev);
		this.placeholder.forget();
		ev.preventDefault();
	}

	_drop(ev){
		var after = this.placeholder.after;
		var data = ev.dataTransfer.getData(ItemType);

		ev.dataTransfer.dropEffect = this.dropEffectFor(ev);
		var dropEffect = ev.dataTransfer.dropEffect;

		if(data){
			// fedwiki type
			var data = JSON.parse(data);
			var item = data.item;

			if(data.url !== this.page.url){
				item.origin = data.url;
				item.originId = item.id;
			}

			// if we make a copy or move to another page
			// we should update the id in the process
			if((dropEffect === 'copy') || (data.url !== this.page.url)){
				item.id = ObjectId();
			}

			// are we moving on the same page?
			if((dropEffect === 'move') && (data.url === this.page.url)){
				// we should not broadcast the moving on the same page
				ev.dataTransfer.dropEffect = 'none';
				this.proxy.modify({
					type: 'move',
					id: item.id,
					after: after
				});
				return;
			}

			this.proxy.modify({
				type: 'add',
				id: item.id,
				item: item,
				after: after
			});

			if(data.created && this.onEditItem){
				this.onEditItem(item.id);
			}
		} else {
			//TODO: always use delayed loading
			var image = dataTransferImage(ev.dataTransfer)
			if(image) {
				var id = ObjectId();
				this.proxy.modify({
					type: "add",
					id: id,
					item: {
						"id": id,
						"type": "image",
						"text": "",
						"url": ""
					},
					after: after
				})

				var self = this;
				var reader = new FileReader();
				reader.onload = function(ev){
					console.log(ev);
					self.proxy.modify({
						type: "edit",
						id: id,
						item: {
							"id": id,
							"type": "image",
							"text": "",
							"url": ev.target.result
						}
					});
				};
				reader.readAsDataURL(image);
				return;
			};

			var item = dataTransferToItem(ev.dataTransfer);
			if(item){
				item.id = ObjectId();
				this.proxy.modify({
					type: "add",
					id: item.id,
					item: item,
					after: after
				});
			}
			return;
		}
	}

	over(ev){
		if(this.page.readonly){
			ev.dataTransfer.dropEffect = 'none';
			return;
		}
		ev.dataTransfer.dropEffect = this.dropEffectFor(ev);

		var size = ev.dataTransfer.getData(SizeType);
		if(size){
			this.placeholder.size = Math.max(parseInt(size), 80);
		} else {
			this.placeholder.size = 30;
		}

		var story = findStory(ev.target);
		if(story === null){
			ev.dataTransfer.dropEffect = 'none';
			this.placeholder.forget();
			return;
		}
		ev.preventDefault();

		var mouse = {
			x: ev.pageX,
			y: ev.pageY
		};

		var originalId = ev.dataTransfer.getData(IdType);
		for(var i = 0; i < story.children.length; i += 1){
			var child = story.children[i];
			if(child.dataset.itemid == originalId){
				continue;
			}

			var box = child.getBoundingClientRect();
			if(mouse.y > box.bottom){
				continue;
			}
			if(mouse.y < (box.top + box.bottom) / 2){
				this.placeholder.moveBefore(child);
			} else {
				this.placeholder.moveAfter(child);
			}
			return;
		}

		this.placeholder.moveAfter(story.lastChild);
	}
}

export class Placeholder {
	constructor(){
		this._after = null;
		this.node = document.createElement("div");
		this.node.classList.add('drop-placeholder');
	}

	set size(value){
		this.node.style.height = value + 'px';
	}

	get after(){
		if(this._after){
			return this._after.dataset.itemid;
		}
		return null;
	}

	forget(){
		this._after = null;
		if(this.node.parentNode){
			this.node.parentNode.removeChild(this.node);
		}
	}

	moveBefore(node){
		if(node == this.node) {
			return;
		}

		this._after = node.previousSibling;
		if(this._after){
			if(this._after.classList.contains("drop-placeholder")){
				this._after = this._after.previousSibling;
			}
		}

		node.parentNode.insertBefore(this.node, node);
	}
	moveAfter(node){
		if(node == this.node) {
			this._after = node.previousSibling;
			return;
		}
		this._after = node;
		if(node){
			node.parentNode.insertBefore(this.node, node.nextSibling);
		}
	}
}


// finds the story and starting node offset relative to node
function findStory(node){
	while(node){
		if(node.classList.contains("story")){
			return node;
		};
		if(node.classList.contains("page")){
			return node.getElementsByClassName("story")[0];
		}
		node = node.parentElement;
	}
	return null;
}

function mouseOffset(ev){
	ev = ev.nativeEvent || ev;
	return {
		x: ev.offsetX || ev.layerX || 0,
		y: ev.offsetY || ev.layerY || 0
	};
}

function dataTransferImage(dataTransfer) {
	var acceptedImages = {"image/png": true, "image/jpeg": true};
	for(var i = 0; i < dataTransfer.files.length; i += 1){
		var file = dataTransfer.files[i];
		if(acceptedImages[file.type]){
			return file;
		}
	}
	return null;
}

function dataTransferToItem(dataTransfer){
	var html = dataTransfer.getData("text/html");
	var href = dataTransfer.getData("text/uri-list");
	if(href){
		if(html){
			var rxTags = /<[^>]+>/g;
			html = html.replace(rxTags, "");
		} else {
			html = URLToTitle(href);
		}
		return {
			id: ObjectId(),
			type: "reference",
			title: html,
			url: href
		};
	}

	// try text
	var text = dataTransfer.getData("text/plain");
	if(text){
		var rxCode = /[=><;{}\[\]]/
		if(text.match(rxCode)){
			return {
				id: ObjectId(),
				type: "code",
				text: text
			};
		}

		return {
			id: ObjectId(),
			type: "paragraph",
			text: text
		};
	}

	console.log("Unhandled drop item:", Object.clone(dataTransfer));
	return null;
}
