package("kb.item.content", function(exports) {
	"use strict";

	depends("DITAIndex.css");

	var Item = React.createClass({
		displayName: "DITAIndexItem",
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

			if (item.children.length > 0) {
				var icon = (expanded ? "mdi-minus" : "mdi-plus");
				if (item.active || item.activechild) {
					icon += "-circle-outline";
				}

				toggle = React.DOM.span({
					className: "dita-index-toggle mdi " + icon,
					onClick: this.toggle
				});

				if (expanded || item.activechild) {
					children = React.DOM.div({
						className: "dita-index-children"
					}, item.children.map(function(item, i) {
						return React.createElement(Item, {
							key: i,
							item: item
						});
					}));
				}
			}

			var link = null;
			if (item.slug === "") {
				link = React.DOM.span({
					className: "dita-index-title " + (item.active ? "dita-index-title-active" : ""),
					onClick: this.open
				}, item.title);
			} else {
				link = React.DOM.a({
					className: "dita-index-title " + (item.active ? "dita-index-title-active" : ""),
					href: item.slug,
					"data-link": item.slug,
					onClick: this.open
				}, item.title);
			}

			return React.DOM.div({
					className: "dita-index-item"
				},
				toggle,
				link,
				children
			);
		}
	});

	// builds a item tree that contains the active/activechild properties
	function build(root, stages) {
		var isactive = {};
		for (var i = 0; i < stages.length; i++) {
			var stage = stages[i];
			if (stage.page.slug !== "") {
				isactive[stage.page.slug] = true;
			}
		}

		function mknode(item) {
			var n = {
				title: item.title,
				slug: item.slug,
				active: isactive[item.slug] || isactive["/" + item.slug],
				activechild: false,
				children: []
			};

			if (Array.isArray(item.children)) {
				for (var i = 0; i < item.children.length; i++) {
					var c = mknode(item.children[i]);
					n.activechild = n.activechild || c.active || c.activechild;
					n.children.push(c);
				}
			}
			return n;
		}

		if (typeof root === "undefined" || root === null) {
			return {
				title: "",
				slug: "",
				active: false,
				activechild: false,
				children: []
			};
		}
		return mknode(root);
	}

	exports["dita-index"] = React.createClass({
		displayName: "DITAIndex",
		contextTypes: {
			Lineup: React.PropTypes.object
		},
		getInitialState: function() {
			return {
				root: build(this.props.item.root, this.context.Lineup.stages)
			};
		},
		activeChanged: function() {
			this.setState({
				root: build(this.props.item.root, this.context.Lineup.stages)
			});
		},
		componentDidMount: function() {
			this.context.Lineup.on("changed", this.activeChanged, this);
		},
		componentWillUnmount: function() {
			this.context.Lineup.remove(this);
		},
		render: function() {
			if (this.props.item === null) {
				return React.DOM.div({
					className: "item-content content-dita-index"
				}, "No index available.");
			}

			var root = this.state.root;
			return React.DOM.div({
					className: "item-content content-dita-index"
				},
				root.children.map(function(item, i) {
					return React.createElement(Item, {
						key: i,
						item: item
					});
				})

			);
		}
	});
});
