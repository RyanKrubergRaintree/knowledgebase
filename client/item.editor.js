import "convert.js";
import "item.js";

'use strict';

KB.Item.Editor = React.createClass({
	commit: function(){
		var stage = this.props.stage,
			item = this.props.item;

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
			var next = Object.clone(item);
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

		if(ev.keyCode == 13){
			switch(item.type){
			case "paragraph":
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

				var adding = {
					type: "paragraph",
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
		}
	},
	render: function(){
		var stage = this.props.stage,
			item = this.props.item;
		return React.DOM.div(
			{ className: "item-content content-editor" },
			React.DOM.textarea({
				ref: "text",
				defaultValue: item.text,
				onBlur: this.commit,
				onKeyDown: this.handleKey,
				autoFocus: true
			})
		)
	}
});
