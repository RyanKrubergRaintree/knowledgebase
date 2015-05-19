'use strict';

import React from 'react';

import {convert} from 'wiki/convert';
import {newIcon} from 'util/icon';
import {TagsInput} from 'util/TagsInput';

var Stats = React.createClass({
	displayName: 'Stats',
	render: function(){
		var self = this;
		return React.DOM.div({className: 'page-stats'},
			React.DOM.div({
				className: 'page-stat helpfulness',
				title: 'Helpfulness (+1/-1 with L/R mouse button)',
				onClick: function(ev){
					self.props.proxy.modify({type: 'vote-up'});
					ev.preventDefault();
				},
				onContextMenu: function(ev){
					self.props.proxy.modify({type: 'vote-down'});
					ev.preventDefault();
				}
			}, 'H:', this.props.rank)
		);
	}
});

var Tags = React.createClass({
	displayName: 'Tags',
	render: function(){
		var tags = this.props.tags.map(function(tag){
			return React.DOM.a({
				key: tag,
				href: '/index/tag/' + tag,
				title: 'Tag: ' + tag,
				className:'tag'
			}, tag);
		});
		return React.DOM.div({className:'page-tags'}, tags)
	}
});

export var Header = React.createClass({
	displayName: 'Header',
	render: function(){
		var proxy = this.props.proxy,
			page = proxy.page;

		var citations = convert.URLToLocation(page.url);
		citations.pathname = '/index/citations' + citations.pathname;

		var pencil = newIcon('pencil', {
			onClick: this.props.onEdit,
			title: 'Edit page header.'
		});

		var meta = React.DOM.div({className: 'page-meta'},
			React.createElement(Tags, {tags: page.meta.tags}),
			React.createElement(Stats, {
				rank: page.meta.upvotes - page.meta.downvotes,
				proxy: this.props.proxy
			})
		);

		if(page.readonly && !page.editableheader){
			pencil = null;
		}
		if(page.dynamic){
			meta = null;
		}

		return React.DOM.div({
			className: 'page-header'
		}, pencil,
		newIcon('wifi', {
			isButton: true,
			title: 'Index of pages citing this page.',
			href: citations.url
		}),
		React.DOM.h1({}, page.title),
		meta,
		React.DOM.div({
			className: 'page-synopsis'
		}, page.synopsis));
	}
});

export var HeaderCreating = React.createClass({
	displayName: 'Header',
	getInitialState: function(){
		return {link: this.props.initialLink};
	},
	update: function(){
		var link = this.refs.link.getDOMNode().value.trim();
		this.setState({link: link});
	},
	create: function(ev){
		var link = this.refs.link.getDOMNode().value,
			tags = this.refs.tags.getTags(),
			synopsis = this.refs.synopsis.getDOMNode().value;

		if(!link){
			return;
		}

		this.props.proxy.create({
			link: link,
			tags: tags,
			synopsis: synopsis
		});
		this.props.onStopEditing();
		ev.preventDefault();
	},
	render: function(){
		var link = this.state.link;
		return React.DOM.div({ className: 'page-header' },
			React.DOM.h1({}, convert.LinkToTitle(link) || '(Untitled)'),
			React.DOM.label({}, 'Link:'),
			React.DOM.input({
				ref: 'link',
				onChange: this.update,
				autoFocus: true,
				defaultValue: link
			}),

			React.DOM.label({}, 'URL:'),
			convert.LinkToURL(link),

			React.DOM.label({}, 'Tags:'),
			React.createElement(TagsInput, {
				ref: 'tags',
				defaultValue: []
			}),

			React.DOM.label({}, 'Synopsis:'),
			React.DOM.textarea({
				className: 'page-synopsis',
				ref: 'synopsis'
			}),

			React.DOM.button({
				onClick: this.create
			}, 'Create New Page')
		);
	}
});

export var HeaderEditing = React.createClass({
	displayName: 'Header',

	submit: function(ev){
		var tags = this.refs.tags.getTags(),
			synopsis = this.refs.synopsis.getDOMNode().value.trim();

		this.props.proxy.modify({
			type: 'header',
			tags: tags,
			synopsis: synopsis
		});
		this.props.onStopEditing();
		ev.preventDefault();
	},
	cancel: function(ev){
		this.props.onStopEditing();
		ev.preventDefault();
	},
	delete: function(ev){
		var link = window.prompt("Please enter page link to confirm deletion.");
		if(link == ""){
			return;
		}
		var verify = convert.URLToLocation(LinkToURL(link||""));
		var actual = convert.URLToLocation(this.page.url);
		if(verify.pathname != actual.pathname){
			return;
		}

		this.props.proxy.delete();
		this.props.onStopEditing();
		ev.preventDefault();
	},

	render: function(){
		var page = this.props.proxy.page;

		return React.DOM.div({className: 'page-header'},
			React.DOM.h1({}, page.title),

			React.DOM.label({}, 'URL:'), page.url,

			React.DOM.label({}, 'Tags:'),
			React.createElement(TagsInput, {
				ref: 'tags',
				autoFocus: true,
				defaultValue: page.meta.tags
			}),

			React.DOM.label({}, 'Synopsis:'),
			React.DOM.textarea({
				ref: 'synopsis',
				className: 'page-synopsis',
				defaultValue: page.synopsis
			}),

			React.DOM.button({ onClick: this.submit }, 'Update'),
			React.DOM.button({ onClick: this.cancel }, 'Cancel'),

			iff(!page.readonly, React.DOM.button({
				style: {'float': 'right'},
				onClick: this.delete
			}, 'Delete Page'))
		);
	}
});
