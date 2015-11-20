package("kb.boot", function(exports) {
	"use strict";

	depends("Auth.js");
	depends("Login.css");

	var LoginForm = React.createClass({
		login: function(ev) {
			ev.preventDefault();
			ev.stopPropagation();

			this.props.provider.login(
				this.refs.username.value,
				this.refs.password.value
			);
		},
		render: function() {
			return React.DOM.form({
					className: "logins",
					onSubmit: this.login
				},
				React.DOM.table(null, React.DOM.tbody(null,
					React.DOM.tr(null,
						React.DOM.td(null, React.DOM.label({
							htmlFor: "username"
						}, "Username:")),
						React.DOM.td(null, React.DOM.input({
							ref: "username",
							name: "username",
							tabIndex: 1
						})),
						React.DOM.td({
								rowSpan: 2
							},
							React.DOM.input({
								className: "button",
								type: "submit",
								value: "Login",
								tabIndex: 3
							}))
					), React.DOM.tr(null,
						React.DOM.td(null, React.DOM.label({
							htmlFor: "password"
						}, "Password:")),
						React.DOM.td(null, React.DOM.input({
							ref: "password",
							name: "password",
							type: "password",
							tabIndex: 2
						}))
					)
				))
			);
		}
	});

	var LoginButton = React.createClass({
		click: function() {
			this.props.provider.login();
		},
		render: function() {
			return React.DOM.div({
				className: "button",
				onClick: this.click
			}, this.props.provider.title);
		}
	});

	var loginView = {
		"form": LoginForm,
		"button": LoginButton
	};

	//TODO: move this logic to server conf
	var loginTitle = {
		"guest": "Customer Login:",
		"google": "Employee Login:"
	};

	exports.Login = React.createClass({
		getInitialState: function() {
			return {
				authLoaded: false,
				error: this.props.initialError
			};
		},
		loaded: function() {
			this.setState({
				authLoaded: true
			});
		},
		componentDidMount: function() {
			if (kb.Auth.loaded) {
				this.loaded();
			} else {
				kb.Auth.on("loaded", this.loaded, this);
			}

			kb.Auth.on("login-error", this.updateError, this);
		},
		componentWillUnmount: function() {
			kb.Auth.remove(this);
		},
		updateError: function(event) {
			this.setState({
				error: event.message
			});
		},
		render: function() {
			var failure = null;
			if (this.state.error !== "") {
				failure = new React.DOM.div({
						className: "login-failed"
					},
					React.DOM.h2(null, "Login failed"),
					React.DOM.p(null, this.state.error)
				);
			}

			var logins = [];
			if (this.state.authLoaded) {
				var providers = kb.Auth.providers;
				for (var name in providers) {
					if (!providers.hasOwnProperty(name)) {
						continue;
					}

					var provider = providers[name];
					var view = loginView[provider.view];
					if (typeof view === "undefined" || view === null) {
						continue;
					}

					logins.push(React.DOM.h2({
						key: "header-" + name
					}, loginTitle[name]));

					logins.push(React.createElement(view, {
						key: name,
						name: name,
						provider: provider
					}));
				}
			}

			return React.DOM.div({
					id: "login"
				},
				React.DOM.div({
					id: "header"
				}, "Knowledge Base"),
				React.DOM.div({
					id: "content"
				}, React.DOM.div({
						className: "modal"
					},

					failure,
					logins
				))
			);
		}
	});
});
