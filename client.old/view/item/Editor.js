'use strict';

import React from 'react';

export var Editor = React.createClass({
	displayName: 'Editor',

	commit: function(){
		var proxy = this.props.proxy;
		var value = this.refs.text.getDOMNode().value;
		var item = this.props.item;
		if((value === '') &&
			((item.type === 'paragraph') || (item.type === 'html'))){
			this.delete();
			return;
		}
		if(value != item.text){
			item.text = value;
			proxy.modify({
				type: 'edit',
				id: item.id,
				item: item
			});
		}
		this.cancel();
	},
	cancel: function(){
		this.props.onCancelEdit();
	},
	delete: function(){
		this.props.proxy.modify({
			type: 'remove',
			id: this.props.item.id
		});
		this.cancel();
	},
	keypress: function(ev){
		if(ev.which == 27 /* escape */){
			this.cancel();
			ev.preventDefault();
		}

		var item = this.props.item,
			page = this.props.page,
			node = this.refs.text.getDOMNode(),
			proxy = this.props.proxy;

		/*
		MERGING of items
		if(ev.which == 8){
			switch(item.type){
			case "paragraph":
				// merge with previous if at the start
				if(node.selectionStart > 0){ return; }

				var i = page.indexOf(item.id);
				// are we the first item?
				if(i <= 0){ return; }

				// is the previous item of the same type?
				var prev = page.story[i1];
				if(prev.type != item.type){
					return;
				}
				var prev = Object.clone(prev);

				// concat the text from two items
				prev.text += node.value;

				// commit the editing of previous item
				if(node.value != ""){
					this.props.onAddToJournal({
						type: "edit",
						id: prev.id,
						item: prev
					});
				}

				// remove the current item
				this.props.onAddToJournal({
					type: "remove",
					id: item.id
				});

				ev.preventDefault();
				this.props.onCancelEdit();

				StartEditingNext = prev.id;

				break;
			}
		}*/

		if(ev.which == 13){
			switch(item.type){
			case "paragraph":
				var pre = node.value.substr(0, node.selectionStart),
					post = node.value.substr(node.selectionStart);

				if(pre != node.value){
					node.value = pre;
				}
				this.commit();

				var adding = {
					type: "paragraph",
					id: ObjectId(),
					text: post
				};

				proxy.modify({
					type: 'add',
					id: adding.id,
					after: item.id,
					item: adding
				});

				this.commit();
				ev.preventDefault();
				break;
			}
		}
	},

	render: function() {
		var item = this.props.item;

		return React.DOM.div({ className: this.props.className},
			React.DOM.div({ className: 'item-edit-toolbar' },
				React.DOM.label({}, item.type, " "),
				React.DOM.button({onClick: this.commit}, 'Save'),
				React.DOM.button({onClick: this.cancel}, 'Cancel'),
				React.DOM.button({onClick: this.delete}, 'Delete')
			),
			React.DOM.textarea({
				ref: 'text',
				onKeyDown: this.keypress,
				defaultValue: item.text,
				autoFocus: true
			})
		);
	}
});
