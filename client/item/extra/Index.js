package("kb.item.content", function(exports) {
	"use strict";

	depends("Index.css");

	var Item = React.createClass({
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
			if (!item.visible) {
				return null;
			}
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

			var className = "index-title " +
				(item.active ? "index-title-active " : "") +
				(item.filter ? "index-title-filter" : "");

			var link = null;
			if (item.slug === "") {
				link = React.DOM.span({
					className: className,
					onClick: this.open
				}, item.title);
			} else {
				link = React.DOM.a({
					className: className,
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
	function build(root, stages, filter) {
		var isactive = {};
		for (var i = 0; i < stages.length; i++) {
			var stage = stages[i];
			if (stage.page.slug !== "") {
				isactive[stage.page.slug] = true;
			}
		}

		filter = filter.toLowerCase();

		function mknode(item) {
			var n = {
				title: item.title,
				slug: item.slug,
				visible: (filter === ""),
				visiblechild: false,
				filter: false,
				filterchild: false,
				active: isactive[item.slug] || isactive["/" + item.slug],
				activechild: false,
				children: []
			};
			n.visible = n.visible || n.active;
			n.visiblechild = n.visiblechild || n.activechild;

			if (Array.isArray(item.children)) {
				for (var i = 0; i < item.children.length; i++) {
					var c = mknode(item.children[i]);
					n.activechild = n.activechild || c.active || c.activechild;
					n.filterchild = n.filterchild || c.filter || c.filterchild;
					n.visiblechild = n.visiblechild || c.visible || c.visiblechild;
					n.children.push(c);
				}
			}

			if (!n.active && !n.activechild && !n.filterchild && (filter !== "")) {
				n.filter = item.title.toLowerCase().indexOf(filter) >= 0;
			}
			n.activechild = n.activechild || n.filter || n.filterchild;
			n.visible = n.visible || n.visiblechild || n.activechild;

			return n;
		}

		if (typeof root === "undefined" || root === null) {
			return {
				title: "",
				slug: "",
				active: false,
				activechild: false,
				visible: true,
				visiblechild: false,
				children: []
			};
		}
		return mknode(root);
	}

	exports["index"] = React.createClass({
		displayName: "Index",
		contextTypes: {
			Lineup: React.PropTypes.object
		},
		getInitialState: function() {
			return {
				filter: "",
				root: build(this.props.item.root, this.context.Lineup.stages, "")
			};
		},
		activeChanged: function() {
			this.setState({
				root: build(this.props.item.root, this.context.Lineup.stages, this.state.filter)
			});
		},
		componentDidMount: function() {
			this.context.Lineup.on("changed", this.activeChanged, this);
		},
		componentWillUnmount: function() {
			window.clearTimeout(this.filterUpdateDelay);
			this.context.Lineup.remove(this);
		},

		filterUpdateDelay: null,
		nextFilter: "",
		changeFilter: function() {
			this.nextFilter = this.refs.filter.value;
			window.clearTimeout(this.filterUpdateDelay);

			var self = this;
			this.filterUpdateDelay = window.setTimeout(function() {
				self.setState({
					filter: self.nextFilter,
					root: build(self.props.item.root, self.context.Lineup.stages, self.nextFilter)
				});
			}, 500);
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
				React.DOM.input({
					ref: "filter",
					className: "index-filter",
					placeholder: "Filter",
					defaultValue: this.state.filter,
					onInput: this.changeFilter
				}),
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
