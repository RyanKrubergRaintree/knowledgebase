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

	componentWillMount: function() {
		this.proxy = this.props.Store.proxyFor(this.props.pageref);
		this.cancel = this.proxy.listen(this.changed);

		this.proxy.initIfNeeded();
		this.proxy.reload();
	},
	componentDidMount: function() {
		this.scrollIntoView();
	},
	componentWillUnmount: function() {
		this.cancel();
		this.proxy = null;
		this.cancel = null;
		this.dragContext = null;
		this.dropArea = null;
	},

	getInitialState: function(){
		return {
			page: null,
			wide: false,
			comments: false,
			mode: 'loading',
			status: {
				state: '',
				code: 0,
				text: '',
				response: ''
			}
		};
	},

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


	triggerEditFor: [],
	editItem: function(id){
		this.triggerEditFor.push(id);
		this.forceUpdate();
	},

	changed: function(proxy){
		var mode = proxy.state;
		if(mode == 'not-found'){
			mode = 'creating';
		}
		this.setState({
			mode: mode,
			page: proxy.page,
			status: {
				state: proxy.state,
				code: proxy.status,
				text: proxy.statusText,
				response: proxy.responseText
			}
		});
	},

	close: function(){ this.props.Lineup.close(this.props.pageref.key); },
	toggleWidth: function(){ this.setState({wide: !this.state.wide}); },
	toggleComments: function(){ this.setState({comments: !this.state.comments}); },

	startEditing: function(){ this.setState({ mode: 'editing' }); },
	stopEditing: function(){ this.setState({ mode: 'loaded' }); },

	scrollIntoView: function(ev){
		if(typeof ev == 'undefined'){
			SmoothScroll.to(this.getDOMNode());
		} else if (!ev.defaultPrevented){
			SmoothScroll.to(this.getDOMNode());
		}
	},

	render: function(){
		var self = this,
			proxy = this.proxy,
			page = this.state.page,
			pageref = this.props.pageref,
			status = this.state.status;

		if((page === null) || (status.state === '')) {
			return React.DOM.article({className: 'page'});
		}

		var props = {
			className: 'page',
			onClick: this.scrollIntoView,
			'data-key': pageref.key,
			'data-url': pageref.url
		};

		if(this.state.wide){
			props.className += ' page-wide';
		}

		var statusbar = React.createElement(StatusBar, {
			className: 'page-top',

			page: page,
			status: this.state.status,
			wide: this.state.wide,

			onToggleWidth: this.toggleWidth,
			onClose: this.close
		});

		var header = null;
		//TODO: push mode into Header
		switch(this.state.mode){
		default:
			header = React.createElement(Header, {
				proxy: proxy,
				onEdit: this.startEditing
			});
			break;
		case 'creating':
			header = React.createElement(HeaderCreating, {
				initialLink: pageref.link,
				proxy: proxy,
				onStopEditing: this.stopEditing
			});
			break;
		case 'editing':
			header = React.createElement(HeaderEditing, {
				proxy: proxy,
				onStopEditing: this.stopEditing
			});
			break;
		}

		if((this.state.mode === 'creating') || (this.state.status.state === 'loading')){
			return React.DOM.article(props,
				statusbar,
				React.DOM.div({
					className: 'page-middle',
					ref: 'container'
				}, header),
				React.DOM.div({className: 'page-bottom'})
			);
		}

		var comments = null;
		if(this.state.comments){
			comments = React.createElement(Comments, {
				proxy: proxy,
				comments: page.comments || []
			});
		}

		var triggerEditFor = this.triggerEditFor;
		this.triggerEditFor = [];

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
