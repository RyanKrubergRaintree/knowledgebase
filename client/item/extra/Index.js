package("kb.item.content", function(exports) {
	"use strict";

	depends("Index.css");

	var Item = createReactClass({
		displayName: "IndexItem",
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
					className: "index-toggle mdi " + icon,
					onClick: this.toggle
				});

				if (expanded || item.activechild) {
					children = React.DOM.div({
						className: "index-children"
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
					className: "index-title " + (item.active ? "index-title-active" : ""),
					onClick: this.open
				}, item.title);
			} else {
				link = React.DOM.a({
					className: "index-title " + (item.active ? "index-title-active" : ""),
					href: item.slug,
					"data-link": item.slug,
					onClick: this.open
				}, item.title);
			}

			return React.DOM.div({
					className: "index-item"
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

	exports["index"] = createReactClass({
		displayName: "Index",
		contextTypes: {
			Lineup: kb.react.object
		},
		getInitialState: function() {
			return {
				root: build(this.props.item.root, this.context.Lineup.stages)
			};
		},
		debounce: 0,
		activeChanged: function() {
			window.clearTimeout(this.debounce);
			this.debounce = window.setTimeout(this.rebuildTree, 100);
		},
		rebuildTree: function() {
			window.clearTimeout(this.debounce);
			this.setState({
				root: build(this.props.item.root, this.context.Lineup.stages)
			});
		},
		componentDidMount: function() {
			this.context.Lineup.on("changed", this.activeChanged, this);
		},
		componentWillUnmount: function() {
			window.clearTimeout(this.debounce);
			this.context.Lineup.remove(this);
		},
		render: function() {
			if (this.props.item === null) {
				return React.DOM.div({
					className: "item-content content-index"
				}, "No index available.");
			}

			var root = this.state.root;
			return React.DOM.div({
					className: "item-content content-index"
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

	// backwards compatibility
	exports["dita-index"] = exports["index"];
});
