'use strict'

import React from 'react';

import {newIcon} from 'util/icon';

var Comment = React.createClass({
	displayName: 'Comment',
	deleteComment: function(){
		this.props.onDelete(this.props.comment.id);
	},
	render: function(){
		var comment = this.props.comment,
			time = new Date(comment.time);
		return React.DOM.div({className: 'comment'},
			newIcon('delete', {
				className: 'icon-red comment-delete',
				onClick: this.deleteComment
			}),
			React.DOM.span({className: 'time'}, time.toLocaleString()),
			React.DOM.p({}, comment.text)
		);
	}
});

export var Comments = React.createClass({
	displayName: 'Comments',
	addComment: function(text){
		var id = ObjectId();
		this.props.proxy.modify({
			type: 'comment-add',
			id: id,
			comment: {
				id: id,
				text: text,
				time: (new Date()).toISOString()
			}
		});
	},
	deleteComment: function(id){
		this.props.proxy.modify({
			type: 'comment-remove',
			id: id
		});
	},
	submit: function(ev){
		var node = this.refs.text.getDOMNode();
		var text = node.value;
		if(text !== ''){
			node.value = '';
			this.addComment(text);
		}
		ev.preventDefault();
	},
	render: function(){
		var self = this;
		var comments = this.props.comments.map(function(c){
			return React.createElement(Comment, {
				key: c.id,
				comment: c,
				onDelete: self.deleteComment
			});
		});
		return React.DOM.div({ className: 'page-comments' },
			React.DOM.div({className: 'comments' }, comments),
			React.DOM.form({onSubmit: this.submit},
				React.DOM.input({
					ref: 'text',
					autoFocus: true
				})
			)
		);
	}
});
