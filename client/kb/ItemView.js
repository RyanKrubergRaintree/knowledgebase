// import /kb/Convert.js

ItemView = {};

'use strict';

ItemView.Unknown = React.createClass({
	displayName: 'Unknown',
	render: function(){
		return React.DOM.div({
			className: 'item-content content-unknown'
		}, this.props.item.text);
	}
});

ItemView['image'] = React.createClass({
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

ItemView['paragraph'] = React.createClass({
	displayName: 'Paragraph',
	render: function(){
		var stage = this.props.stage;
		return React.DOM.p({
			className: 'item-content content-paragraph',
			dangerouslySetInnerHTML: {
				__html: stage.resolveLinks(this.props.item.text)
			}
		});
	}
});

ItemView['html'] = React.createClass({
	displayName: 'HTML',
	render: function(){
		var stage = this.props.stage;
		return React.DOM.div({
			className: 'item-content content-html',
			dangerouslySetInnerHTML: {
				__html: stage.resolveLinks(this.props.item.text)
			}
		});
	}
});

ItemView['code'] = React.createClass({
	displayName: 'Code',
	render: function(){
		return React.DOM.div({
			className: 'item-content content-code'
		}, this.props.item.text);
	}
});

ItemView['reference'] = React.createClass({
	displayName: 'Reference',
	render: function(){
		var item = this.props.item;
		var url = Convert.URLToLocation(item.url);
		var external = url.origin && (url.origin != window.location.origin);

		//TODO: handle external links
		var url = item.url || Convert.LinkToURL(item.title);
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

ItemView['entry'] = React.createClass({
	displayName: 'Entry',
	render: function(){
		var item = this.props.item;
		var url = Convert.URLToLocation(item.url);
		var external = url.origin && (url.origin != window.location.origin);

		var origin = url.origin || window.location.origin;

		var className = external ? 'external-link': '';
		var target = external ? '_blank': '';

		var url = item.url || Convert.LinkToURL(item.title);

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
