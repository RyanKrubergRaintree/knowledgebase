package('kb.item.content', function(exports) {
	'use strict';

	depends('DITAIndex.css');

	var Item = React.createClass({
		displayName: 'DITAIndexItem',
		getInitialState: function() {
			return {
				expanded: false
			};
		},
		open: function() {
			this.setState({
				expanded: true
			});
		},
		toggle: function() {
			this.setState({
				expanded: !this.state.expanded
			});
		},
		render: function() {
			var item = this.props.item;

			var expanded = this.state.expanded;

			var toggle = null;
			var children = null;

			if (Array.isArray(item.children)) {
				toggle = React.DOM.span({
					className: 'dita-index-toggle mdi ' + (expanded ? 'mdi-minus' : 'mdi-plus'),
					onClick: this.toggle
				});

				if (expanded) {
					var opened = this.props.opened;
					children = React.DOM.div({
						className: 'dita-index-children'
					}, item.children.map(function(item, i) {
						return React.createElement(Item, {
							key: i,
							item: item,
							opened: opened
						});
					}));
				}
			}

			var isactive = this.props.opened.indexOf(item.slug) >= 0;

			var link = null;
			if (item.slug === '') {
				link = React.DOM.span({
					className: 'dita-index-title ' + (isactive ? 'dita-index-title-active' : ''),
					onClick: this.open
				}, item.title);
			} else {
				link = React.DOM.a({
					className: 'dita-index-title ' + (isactive ? 'dita-index-title-active' : ''),
					href: item.slug,
					onClick: this.open
				}, item.title);
			}

			return React.DOM.div({
					className: 'dita-index-item'
				},
				toggle,
				link,
				children
			);
		}
	});

	function extractSlugs(stages) {
		var r = [];
		for (var i = 0; i < stages.length; i++) {
			var stage = stages[i];
			if (stage.page.slug !== '') {
				r.push(stage.page.slug);
			}
		}
		return r;
	}

	exports['dita-index'] = React.createClass({
		displayName: 'DITAIndex',
		getInitialState: function() {
			return {
				opened: []
			};
		},
		activeChanged: function(ev) {
			this.setState({
				opened: extractSlugs(ev.lineup.stages)
			});
		},
		componentDidMount: function() {
			kb.app.Lineup.on('changed', this.activeChanged, this);
		},
		componentWillUnmount: function() {
			kb.app.Lineup.remove(this);
		},
		render: function() {
			if (this.props.item === null) {
				return React.DOM.div({
					className: 'item-content content-dita-index'
				}, 'No index available.');
			}

			var root = this.props.item.root;
			var opened = this.state.opened;
			return React.DOM.div({
					className: 'item-content content-dita-index'
				},
				root.children.map(function(item, i) {
					return React.createElement(Item, {
						key: i,
						item: item,
						opened: opened
					});
				})

			);
		}
	});
});
