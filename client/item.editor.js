import "convert.js";
import "item.js";

'use strict';

KB.Item.Editor = React.createClass({
	blur: function(ev){
		var self = this;
		var container = this.getDOMNode();
		if(ev.relatedTarget === null){
			// HACK-FIX: Firefox doesn't have relatedTarget
			// use a delayed commit if we don't know the related target
			// that way any menu action can complete, before the commit
			window.setTimeout(function(){ self.commit(); }, 300);
			return;
		}
		if(!container.contains(ev.relatedTarget)){
			this.commit();
		}
	},
	commit: function(ev){
		var stage = this.props.stage,
			item = this.props.item;

		// in case the delay
		if(this.refs.text === undefined) { return; }

		var text = this.refs.text.getDOMNode().value;
		if((text === '') &&
			((item.type === 'paragraph') || (item.type === 'html'))){
			stage.patch({
				id: item.id,
				type: "remove"
			});
			return;
		}

		if(item.text != text){
			var next = DeepClone(item);
			next.text = text;
			stage.patch({
				id: next.id,
				type: "edit",
				item: next
			})
		}
		this.stopEditing();
	},
	stopEditing: function(){
		var stage = this.props.stage,
			item = this.props.item;

		try {
			var actual = stage.page.itemById(item.id);
			if(actual && (actual.type == "paragraph") && (actual.text == "")){
				stage.patch({
					id: actual.id,
					type: "remove"
				});
			}
		} catch(ex){}


		stage.editing.stop(item.id);
	},

	remove: function(ev){
		var stage = this.props.stage,
			item = this.props.item;

		stage.patch({
			type: 'remove',
			id: item.id
		});

		ev.preventDefault();
		ev.stopPropagation();
	},

	handleKey: function(ev){
		if(ev.keyCode == 27){
			this.stopEditing();
			ev.preventDefault();
			ev.stopPropagation();
			return;
		}

		var stage = this.props.stage,
			item = this.props.item,
			node = this.refs.text.getDOMNode();

		if((ev.ctrlKey) && (ev.keyCode == 13)){

			var pre = node.value.substr(0, node.selectionStart),
				post = node.value.substr(node.selectionStart);

			if(pre == ""){
				ev.preventDefault();
				return;
			}

			if(pre != node.value){
				node.value = pre;
				this.commit();
				this.stopEditing();
			}

			var continueWith = {
				"paragraph": "paragraph",
				"html": "html"
			};

			var adding = {
				type: continueWith[item.type] || "paragraph",
				id: GenerateID(),
				text: post
			};

			stage.patch({
				type: 'add',
				id: adding.id,
				after: item.id,
				item: adding
			});

			stage.editing.start(adding.id);
			ev.preventDefault();
		}
	},
	render: function(){
		var stage = this.props.stage,
			item = this.props.item;
		return React.DOM.div(
			{
				className: "item-content content-editor",
				onBlur: this.blur
			},
			React.DOM.div({className: "item-type"}, item.type),
			React.DOM.textarea({
				ref: "text",
				defaultValue: item.text,
				onKeyDown: this.handleKey,
				autoFocus: true
			}),
			React.DOM.div(
				{className:'editor-buttons'},
				React.DOM.div({
					className: "mdi mdi-content-save",
					title: "Save changes. (Unfocus the editor will save.)",
					tabIndex: '1',
					onClick: this.commit
				}),
				React.DOM.div({
					className: "mdi mdi-backup-restore",
					title: "Cancel any modifications. (Pressing escape will cancel changes.)",
					tabIndex: '2',
					onClick: this.stopEditing
				}),
				React.DOM.div({
					className: "item-delete mdi mdi-delete",
					title: "Delete this item.",
					tabIndex: '3',
					onClick: this.remove
				})
			)
		)
	}
});
