import "util/SmoothScroll.js";
import "page.js";
import "item.view.js";
import "util/drag.js"
import "drop.js"

KB.Page.View = (function(){
	// ensure that we clear all the place-holders
	document.removeEventListener("dragend", clearDropPosition);
	document.removeEventListener("drop", clearDropPosition);
	document.addEventListener("dragend", clearDropPosition);
	document.addEventListener("drop", clearDropPosition);

	function clearDropPosition(){
		var els = document.getElementsByClassName("drop-after");
		for(var i = 0; i < els.length; i += 1){
			els[i].classList.remove("drop-after");
		}
		var els = document.getElementsByClassName("drop-before");
		for(var i = 0; i < els.length; i += 1){
			els[i].classList.remove("drop-before");
		}
	}

	function getAfter(page, drop){
		if(drop == null){
			return;
		}
		var story = page.story;

		if(drop.node == null){
			if(drop.rel == 'after'){
				if(story.length > 0){
					return story[story.length-1].id; // as last
				}
			}
			return; // as first
		}

		var id = drop.node.dataset.id;
		if(id == null){
			return; // as first
		}
		for(var i = 0; i < story.length; i += 1){
			if(story[i].id == id){
				if(drop.rel == 'after'){
					return id;
				} else if (i > 0) {
					// as first
					return story[i-1].id;
				} else {
					return;
				}
			}
		}
		return;
	}

	var Page = React.createClass({
		displayName: "Page",

		dropEffectFor: function(ev){
			if(ev.dataTransfer.effectAllowed === 'copy'){
				return 'copy';
			}
			if(ev.shiftKey){
				return 'copy';
			}
			return 'move';
		},

		dragEnter: function(ev){},
		dragOver: function(ev){
			clearDropPosition();
			if(!this.props.stage.canModify){
				return;
			}

			var drop = FindListPosition(ev, "page", "page-story");
			if(drop == null){
				return
			}

			ev.preventDefault();
			ev.dataTransfer.dropEffect = this.dropEffectFor(ev);

			var drop = FindListPosition(ev, "page", "page-story");
			if(drop == null){
				return
			}

			if(drop.node != null){
				drop.node.classList.add("drop-" + drop.rel);
			} else {
				var story = this.refs.story.getDOMNode();;
				story.classList.add("drop-" + drop.rel);
			}
		},
		dragDrop: function(ev){
			ev.preventDefault();
			ev.stopPropagation();

			var stage = this.props.stage,
				page = stage.page;

			clearDropPosition();
			if(!stage.canModify){
				return;
			}

			var drop = FindListPosition(ev, "page", "page-story");
			if(drop == null){
				return
			}
			var after = getAfter(page, drop);

			var dropEffect = this.dropEffectFor(ev);

			console.log(ev.nativeEvent);

			var data = ev.dataTransfer.getData("kb/item");
			if(data){
				var data = JSON.parse(data);
				var item = data.item;

				if(data.url && (data.url !== page.url)){
					item.origin = data.url;
					item.originId = item.id;
				}

				// if we make a copy or move to another page
				// we should update the id in the process
				if((dropEffect === 'copy') || (data.url !== stage.url)){
					item.id = GenerateID();
				}

				// are we moving on the same page?
				if((dropEffect === 'move') && (data.url === stage.url)){
					window.DropCanceled = true;
					stage.patch({
						type: 'move',
						id: item.id,
						after: after
					});

					ev.dataTransfer.effectAllowed = 'none';
					ev.dataTransfer.dropEffect = 'none';
					return;
				}

				stage.patch({
					type: 'add',
					id: item.id,
					item: item,
					after: after
				});
				ev.dataTransfer.dropEffect = dropEffect;
			} else {
				ev.dataTransfer.dropEffect = 'copy';
				DropData(stage, after, ev.dataTransfer);
			}
		},
		dragLeave: function(ev){
			//TODO: fix glitchy rendering
			clearDropPosition();
		},
		render: function(){
			var stage = this.props.stage,
				page = this.props.page;

			return React.DOM.div(
				{
					className: "page",

					onDragEnter: this.dragEnter,
					onDragOver: this.dragOver,
					onDrop: this.dragDrop,
					onDragLeave: this.dragLeave
				},
				React.DOM.h1(null, page.title),
				React.createElement(Story, {
					ref: "story",

					stage: stage,
					page: page,
					story: page.story
				})
			);
		}
	});

	var Story = React.createClass({
		displayName: "Story",
		render: function(){
			var stage = this.props.stage,
				page = this.props.page,
				story = this.props.story;

			return React.DOM.div(
				{className: "page-story"},
				story.map(function(item, i){
					return React.createElement(KB.Item.View, {
						key: item.id || i,
						stage: stage,
						item: item
					});
				})
			);
		}
	});

	return Page;
})();
