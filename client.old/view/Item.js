'use strict';

import React from 'react';
import ObjectId from 'ObjectId';

import {convert} from 'wiki/convert';

import {Editor} from './item/Editor';
import {Types, UnknownType} from './item/Types';
import {DragContext} from './Dragging';

var Origin = React.createClass({
	displayName: 'Origin',
	render: function(){
		var item = this.props.item;
		if(!item.origin){
			return React.DOM.div({className: 'item-origin hidden' });
		}
		var url = item.origin;
		return React.DOM.a({
			className: 'item-origin',
			href: item.origin
		});
	}
});

export var Item = React.createClass({
	displayName: 'Item',
	getInitialState: function(){
		return { editing: this.props.editing };
	},
	componentWillReceiveProps: function(nextProps) {
		if(nextProps.editing){
			this.setState({editing: true});
		}
	},

	startEditing: function(){ this.setState({editing: true}); },
	cancelEditing: function(){ this.setState({editing: false}); },

	// dragging
	dragContext: null,
	dragStart: function(ev, node, item){
		this.dragContext = new DragContext(this.props.proxy);
		this.dragContext.start(
			ev,
			this.props.item,
			this.refs.item.getDOMNode(),
			false
		);
	},
	drag: function(ev){
		this.dragContext.drag(ev);
	},
	dragEnd: function(ev){
		this.dragContext.end(ev);
		this.dragContext = null;
	},

	render: function(){
		var proxy = this.props.proxy;
		var page = proxy.page;
		var item = this.props.item;

		if(this.state.editing){
			return React.createElement(Editor, {
				className: 'item item-editing',
				item: item,
				proxy: proxy,
				onCancelEdit: this.cancelEditing
			});
		}

		var classes = '';
		if(this.props.focus){
			classes += 'item-hightlight'
		}
		var Content = Types[item.type] || UnknownType;
		return React.DOM.div({
				ref: 'item',
				className: 'item ' + classes,
				onDoubleClick: iff(!page.readonly, this.startEditing),
				'data-itemid': item.id
			},
			React.DOM.div({
				className: 'item-move',
				draggable: true,

				onDragStart: this.dragStart,
				onDrag: this.drag,
				onDragEnd: this.dragEnd
			}),
			React.createElement(Origin, {item: item}),
			React.createElement(Content, {
				item: item,
				proxy: proxy
			})
		);
	}
});
