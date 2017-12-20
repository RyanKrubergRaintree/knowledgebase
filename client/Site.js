package("kb", function(exports) {
	"use strict";

	depends("site.css");
	depends("Convert.js");
	depends("Lineup.View.js");

	var HeaderMenu = createReactClass({
		displayName: "HeaderMenu",
		render: function() {
			return React.DOM.div({
					className: "header-menu"
				},
				this.props.items.map(function(item) {
					if (item == null) {
						return null;
					}
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

	var Search = createReactClass({
		lastStageId: undefined,
		getInitialState: function() {
			return {
				options: [
					"All",
					"10.2.700",
					"10.2.500",
					"10.2.400"
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
			var defaultFilter = this.props.Session.filter;
			var options = filter.options.slice();
			if (options.indexOf(defaultFilter) < 0) {
				options.push(defaultFilter);
			}
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
						defaultValue: defaultFilter
					},
					options.map(function(item) {
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

	var LoginInfo = createReactClass({
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

	var Header = createReactClass({
		openHome: function(ev) {
			ev.preventDefault();
			ev.stopPropagation();

			var session = this.props.Session;

			var pages = [session.home];
			if (session.filter) {
				pages.push("help-" + kb.convert.TextToSlug(session.filter) + "=Index");
			}

			this.props.Lineup.openPages(pages);
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
			var company = this.props.Session.user.company || "Community";
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
					}, (company != "" ? {
						key: "company",
						href: "/group=" + company,
						caption: company
					} : null), {
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

	var Content = createReactClass({
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

	exports.Site = createReactClass({
		displayName: "Site",
		render: function() {
			return React.DOM.div(null,
				React.createElement(Header, this.props),
				React.createElement(Content, this.props)
			);
		}
	});
});
