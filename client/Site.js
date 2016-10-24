package("kb", function(exports) {
	"use strict";

	depends("site.css");
	depends("Lineup.View.js");

	var HeaderMenu = React.createClass({
		displayName: "HeaderMenu",
		render: function() {
			return React.DOM.div({
					className: "header-menu"
				},
				this.props.items.map(function(item) {
					var props = {
						key: item.key,
						href: item.href,
						onClick: item.onClick
					};
					return React.DOM.a(props, item.caption);
				})
			);
		}
	});

	var Search = React.createClass({
		lastStageId: undefined,
		getInitialState: function() {
			return {
				options: [
					"All",
					"10.2.600",
					"10.2.500",
					"10.2.400",
					"10.2.300",
					"10.2.200",
					"10.2.100",
					"9.4"
				],
				hidden: this.props.Session.branch != ""
			};
		},
		updateFilter: function(ev) {
			var filter = this.refs.filter.value;
			if (filter !== "All") {
				this.props.Session.filter = filter;
			} else {
				this.props.Session.filter = "";
			}
		},
		search: function(ev) {
			var Lineup = this.props.Lineup;
			var query = this.refs.query.value.trim();
			if (ev.shiftKey) {
				this.lastStageId = undefined;
			}

			var param = "q=" + query;

			this.lastStageId = Lineup.open({
				url: "/search=search?" + param,
				title: "\"" + query + "\"",
				insteadOf: this.lastStageId
			});
			ev.preventDefault();
		},
		keyDown: function(ev) {
			var Lineup = this.props.Lineup;
			if (ev.keyCode === 27) { // esc
				Lineup.closeLast();
				return;
			}

			if (ev.keyCode === 13) {
				// open page directly
				if (ev.ctrlKey) {
					if (!ev.shiftKey) {
						Lineup.clear();
					}

					var query = this.refs.query.value.trim();
					Lineup.openLink(query);
					ev.preventDefault();
					return;
				}

				this.search(ev);
				ev.preventDefault();
				return;
			}

			var stages = document.querySelectorAll(".stage");
			if (stages.length === 0) {
				return;
			}

			var stage = stages[stages.length - 1];
			var middle = stage.querySelector(".stage-scroll");

			switch (ev.keyCode) {
				case 33: // pageup
					middle.scrollTop -= middle.clientHeight;
					break;
				case 34: // pagedown
					middle.scrollTop += middle.clientHeight;
					break;
			}
		},
		render: function() {
			var filter = this.state;
			return React.DOM.form({
					className: "search",
					onSubmit: this.search
				},
				React.DOM.input({
					ref: "query",
					placeholder: "Search...",
					onKeyDown: this.keyDown
				}),
				React.DOM.select({
						className: "search-filter",
						ref: "filter",
						style: {
							display: filter.hidden ? "none" : null
						},
						onChange: this.updateFilter,
						defaultValue: this.props.Session.filter
					},
					filter.options.map(function(item) {
						return React.DOM.option({
							key: item,
							value: item
						}, item);
					})
				),
				React.DOM.button({
					className: "search-icon mdi mdi-magnify",
					type: "submit",
					tabIndex: -1
				})
			);
		}
	});

	var LoginInfo = React.createClass({
		render: function() {
			var infostyle = {
				style: {
					fontSize: "smaller"
				}
			};

			var username = this.props.username;
			if (username.indexOf("=") >= 0) {
				var tokens = username.split("=");
				return React.DOM.div({
						className: "background-info"
					},
					React.DOM.div(infostyle, "logged in as:"),
					React.DOM.div({}, tokens[1]),
					React.DOM.div(infostyle, "customer:"),
					React.DOM.div({}, tokens[0])
				);
			}

			return React.DOM.div({
					className: "background-info"
				},
				React.DOM.div(infostyle, "logged in as:"),
				React.DOM.div({}, username)
			);
		}
	});

	var Header = React.createClass({
		openHome: function(ev) {
			ev.preventDefault();
			ev.stopPropagation();

			this.props.Lineup.openPages(this.props.Session.home);
		},
		createNewPage: function(ev) {
			ev.preventDefault();
			ev.stopPropagation();

			var lineup = this.props.Lineup;
			lineup.open({
				url: "",
				link: "",
				title: ""
			});
		},
		logout: function(ev) {
			ev.preventDefault();
			ev.stopPropagation();

			this.props.Session.logout();
		},
		displayName: "Header",
		render: function() {
			return React.DOM.div({
					id: "header"
				},
				React.DOM.a({
					className: "button home mdi mdi-home",
					href: "#",
					title: "Home",
					onClick: this.openHome
				}),
				React.createElement(Search, this.props),
				React.createElement(HeaderMenu, {
					items: [{
						key: "new-page",
						href: "#",
						onClick: this.createNewPage,
						caption: "New Page"
					}, {
						key: "index",
						href: "/help=index",
						caption: "Index"
					}, {
						key: "recent-changes",
						href: "/page=recent-changes",
						caption: "Recent Changes"
					}, {
						key: "user",
						href: "/user=current",
						caption: "User"
					}, {
						key: "logout",
						href: "#",
						onClick: this.logout,
						caption: "Logout"
					}]
				})
			);
		}
	});

	var Content = React.createClass({
		displayName: "Content",
		render: function() {
			return React.DOM.div({
					id: "content"
				},
				React.createElement(LoginInfo, {
					username: this.props.Session.user.name
				}),
				React.createElement(kb.Lineup.View, this.props)
			);
		}
	});

	exports.Site = React.createClass({
		displayName: "Site",
		render: function() {
			return React.DOM.div(null,
				React.createElement(Header, this.props),
				React.createElement(Content, this.props)
			);
		}
	});
});
