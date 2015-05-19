'use strict';

import React from 'react';
import {convert, resolve} from 'wiki/convert';

export var Types = {};

export var UnknownType = React.createClass({
	displayName: 'Unknown',
	render: function(){
		return React.DOM.div({
			className: 'item-content content-unknown'
		}, this.props.item.text);
	}
});

Types['image'] = React.createClass({
	displayName: 'Image',
	render: function(){
		return React.DOM.div({
			className: 'item-content content-image'
		},
			React.DOM.img({src: this.props.item.url}),
			React.DOM.p({}, this.props.item.text)
		);
	}
});

Types['paragraph'] = React.createClass({
	displayName: 'Paragraph',
	render: function(){
		return React.DOM.p({
			className: 'item-content content-paragraph',
			dangerouslySetInnerHTML: {
				__html: resolve(this.props.item.text)
			}
		});
	}
});

Types['html'] = React.createClass({
	displayName: 'HTML',
	render: function(){
		return React.DOM.div({
			className: 'item-content content-html',
			dangerouslySetInnerHTML: {
				__html: resolve(this.props.item.text)
			}
		});
	}
});

Types['code'] = React.createClass({
	displayName: 'Code',
	render: function(){
		return React.DOM.div({
			className: 'item-content content-code'
		}, this.props.item.text);
	}
});

Types['reference'] = React.createClass({
	displayName: 'Reference',
	render: function(){
		var item = this.props.item;
		var url = convert.URLToLocation(item.url);
		var external = url.origin && (url.origin != window.location.origin);

		//TODO: handle external links
		var url = item.url || convert.LinkToURL(item.title);
		return React.DOM.div({className: 'item-content content-reference'},
			React.DOM.a({
				className: external ? 'external-link': '',
				target: external ? '_blank': '',
				href: url
			}, item.title),
			React.DOM.p({}, this.props.item.text)
		);
	}
});

Types['entry'] = React.createClass({
	displayName: 'Entry',
	render: function(){
		var item = this.props.item;
		var url = convert.URLToLocation(item.url);
		var external = url.origin && (url.origin != window.location.origin);

		var origin = url.origin || window.location.origin;

		var className = external ? 'external-link': '';
		var target = external ? '_blank': '';

		var url = item.url || convert.LinkToURL(item.title);

		return React.DOM.div({className: 'item-content content-entry'},
			React.DOM.span({className: 'entry-rank'}, item.rank || 0),
			React.DOM.span({className: 'entry-origin'}, origin),
			React.DOM.a({
				className: className + ' entry-title',
				target: target,
				title: url,
				href: url
			}, item.title),
			React.DOM.p({className: 'synopsis'}, this.props.item.text)
		);
	}
});
