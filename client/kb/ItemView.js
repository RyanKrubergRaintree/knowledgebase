// import /kb/Convert.js
// import /kb/Resolve.js

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
				__html: Resolve(stage, this.props.item.text)
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
				__html: ResolveHTML(stage, this.props.item.text)
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
		var url = Convert.LinkToReference(item.link).url;
		var loc = Convert.URLToLocation(url);
		var external = loc.origin && (loc.origin != window.location.origin);

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
		var ref = Convert.LinkToReference(item.link);
		var url = ref.url;
		return React.DOM.div({
			className: 'item-content content-entry',
		},
			React.DOM.a({
				className: 'entry-title',
				title: url,
				href: url
			}, item.title),
			React.DOM.div({className: 'entry-owner'}, ref.owner),
			React.DOM.p({className: 'entry-synopsis'}, this.props.item.text)
		);
	}
});

ItemView['tags'] = React.createClass({
	displayName: 'Tags',
	render: function(){
		var item = this.props.item,
			stage = this.props.stage;
		var tags = item.text.split(",");

		return React.DOM.div({className: 'item-contet content-tags'},
			tags.map(function(tag, i){
				tag = tag.trim();
				return React.DOM.a({
					className: "tag",
					key: i,
					href: '/index/tag/' + Slugify(tag)
				}, tag);
			}),
			React.DOM.div({className:"clear-fix"})
		);
	}
});
