package('kb.item.content', function(exports){
	'use strict';

	depends('dita.css');
	depends('content.css');

	depends('../Convert.js');

	depends('Sanitize.js');
	depends('Resolve.js');

	exports.Unknown = React.createClass({
		displayName: 'Unknown',
		render: function(){
			var item = this.props.item;
			return React.DOM.div(
				{ className: 'item-content content-unknown' },
				React.DOM.span({style: {'float': 'right'}}, item.type),
				React.DOM.p({}, item.text),
				React.DOM.div({className:'clear-fix'})
			);
		}
	});

	var ContentTypes = [
		{name: 'Text', type: 'paragraph', desc: 'simple text paragraph'},
		{name: 'HTML', type: 'html', desc: 'a subset of html for more advanced content'},
		{name: 'Code', type: 'code', desc: 'item especially designed for code'},
		{name: 'Tags', type: 'tags', desc: 'tags for the page'}
	];

	exports.factory = React.createClass({
		displayName: 'Factory',
		convert: function(ev){
			var type = GetDataAttribute(ev.currentTarget, 'type');
			var stage = this.props.stage,
				item = this.props.item;

			stage.patch({
				type: 'edit',
				id: item.id,
				item: {
					type: type,
					id: item.id,
					text: item.text
				}
			});

			stage.editing.start(item.id);
		},

		render: function(){
			var self = this;
			var item = this.props.item;
			return React.DOM.div(
				{ className: 'item-content content-factory'	},
				React.DOM.p({}, item.text || 'Create new '),
				ContentTypes.map(function(item){
					return React.DOM.button(
						{
							key: item.type,
							className: 'factory-item',
							'data-type': item.type,
							title: item.desc,
							onClick: self.convert
						}, item.name);
				})
			);
		}
	});

	exports.image = React.createClass({
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

	exports.paragraph = React.createClass({
		displayName: 'Paragraph',
		render: function(){
			var stage = this.props.stage;
			var resolved = kb.item.Resolve(stage, this.props.item.text);
			var paragraphs = resolved.split('\n\n');
			if(paragraphs.length > 1){
				return React.DOM.div({
					className: 'item-content content-paragraph'
				}, paragraphs.map(function(p, i){
					return React.DOM.p({key: i, dangerouslySetInnerHTML: {__html: kb.item.Sanitize(p)}});
				}));
			} else {
				return React.DOM.p({
					className: 'item-content content-paragraph',
					dangerouslySetInnerHTML: {
						__html: kb.item.Sanitize(paragraphs[0])
					}
				});
			}
		}
	});

	exports.html = React.createClass({
		displayName: 'HTML',
		render: function(){
			var stage = this.props.stage;
			return React.DOM.div({
				className: 'item-content content-html',
				dangerouslySetInnerHTML: {
					__html: kb.item.Sanitize(kb.item.ResolveHTML(stage, this.props.item.text))
				}
			});
		}
	});

	exports.code = React.createClass({
		displayName: 'Code',
		render: function(){
			return React.DOM.div({
				className: 'item-content content-code'
			}, this.props.item.text);
		}
	});

	exports.reference = React.createClass({
		displayName: 'Reference',
		render: function(){
			var item = this.props.item;
			var url = item.url;
			var loc = kb.convert.URLToLocation(url);
			var external = (loc.origin !== '') &&
				(loc.origin !== window.location.origin);

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

	exports.entry = React.createClass({
		displayName: 'Entry',
		render: function(){
			var item = this.props.item;
			var ref = kb.convert.LinkToReference(item.link);
			var url = ref.url;
			return React.DOM.div({
				className: 'item-content content-entry'
			},
				React.DOM.a({
					className: 'entry-title',
					title: url,
					href: url
				}, item.title),
				React.DOM.div({className: 'entry-owner'}, ref.owner),
				React.DOM.p({
					className: 'entry-synopsis',
					dangerouslySetInnerHTML: {
						__html: this.props.item.text
					}
				})
			);
		}
	});

	exports.tags = React.createClass({
		displayName: 'Tags',
		render: function(){
			var item = this.props.item;

			var text = typeof item.text === 'undefined' ? '' : item.text.trim();
			var tags = [];
			if(text !== '') {
				tags =  text.split(',');
			}
			tags = tags.map(function(tag){ return tag.trim(); })
						.filter(function(tag){ return tag !== ''; });
			return React.DOM.div({className: 'item-contet content-tags'},
				tags.length > 0 ?
					tags.map(function(tag, i){
						tag = tag.trim();
						return React.DOM.a({
							className: 'tag',
							key: i,
							href: '/tag:' + kb.convert.TextToSlug(tag)
						}, tag);
					})
				: React.DOM.p({}, 'Double click here to add page tags.'),
				React.DOM.div({className:'clear-fix'})
			);
		}
	});

});
