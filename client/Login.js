package("kb.boot", function(exports) {
	"use strict";

	depends("Session.js");
	depends("Login.css");

	var LoginForm = React.createClass({
		login: function(ev) {
			ev.preventDefault();
			ev.stopPropagation();

			var form = new FormData();
			form.append("user", this.refs.username.value);
			form.append("code", this.refs.password.value);

			kb.Session.fetch({
				url: this.props.url,
				ondone: this.change,
				onerror: this.error,
				body: form
			});
		},
		change: function(response) {
			if (!response.ok) {
				this.props.onFailure(response.text);
				return;
			}
			this.props.onSuccess(response.json);
		},
		error: function(ev) {
			console.log(ev);
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

	var Google = React.createClass({
		render: function() {
			return React.DOM.a({
				className: "button",
				href: "#"
			}, "Google");
		}
	});

	var loginViews = {
		"form": LoginForm,
		"google": Google
	};

	//TODO: move this logic to server conf
	var loginTitle = {
		"guest": "Customer Login:",
		"google": "Employee Login:"
	};

	var order = ["guest", "google"];

	exports.Login = React.createClass({
		getInitialState: function() {
			return {
				error: this.props.initialError
			};
		},
		loginFailed: function(message) {
			this.setState({
				error: message
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

			var self = this;
			var providers = this.props.providers;
			var logins = [];
			order.forEach(function(name) {
				if (!providers.hasOwnProperty(name)) {
					return;
				}

				var params = providers[name];
				var clazz = loginViews[params.kind];
				if (typeof clazz === "undefined" || clazz === null) {
					return;
				}

				logins.push(React.DOM.h2({
					key: "header-" + name
				}, loginTitle[name]));

				logins.push(React.createElement(clazz, {
					key: name,
					url: "/system/auth/" + name,
					params: params,
					onSuccess: self.props.onSuccess,
					onFailure: self.loginFailed
				}));
			});

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
