import React from 'react';

import {StatusBar} from './page/StatusBar';
import {Header, HeaderCreating, HeaderEditing} from './page/Header';
import {Story} from './Story';
import {Comments} from './page/Comments';
import {Footer} from './page/Footer';

import {SmoothScroll} from 'util/SmoothScroll';

import {DropArea} from './Dragging';

//TODO: implement!!!
//  componentWillReceiveProps
//  to handle case when the pageref changes
export var Page = React.createClass({
	displayName: 'Page',

	proxy: null,
	cancel: null,

	dropArea: null,
	dragEnter: function(ev){
		if(this.dropArea === null){
			this.dropArea = new DropArea(this.proxy, this.getAreaContainer, this.editItem);
		}
		this.dropArea.enter(ev);
	},
	dragOver: function(ev){
		this.dropArea.over(ev);
	},
	dragDrop: function(ev){
		this.dropArea.drop(ev);
	},
	dragLeave: function(ev){
		this.dropArea.leave(ev);
	},
	getAreaContainer: function(){
		return this.refs.container.getDOMNode();
	},

	scrollIntoView: function(ev){
		if(typeof ev == 'undefined'){
			SmoothScroll.to(this.getDOMNode());
		} else if (!ev.defaultPrevented){
			SmoothScroll.to(this.getDOMNode());
		}
	},

	render: function(){
		return React.DOM.article(props,
			statusbar,
			React.DOM.div({
				ref: 'container',
				className: 'page-middle',
				onDragEnter: this.dragEnter,
				onDragOver: this.dragOver,
				onDrop: this.dragDrop,
				onDragLeave: this.dragLeave
			},
				header,
				React.createElement(Story, {
					className: 'page-content',
					story: page.story,
					proxy: proxy,
					shouldEdit: triggerEditFor
				})
			),
			comments,
			React.createElement(Footer, {
				className: 'page-bottom',
				proxy: proxy,
				onToggleComments: this.toggleComments
			})
		);
	}
})
