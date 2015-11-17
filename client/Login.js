package("kb.boot", function(exports) {
	"use strict";

	depends("app.css");
	depends("Login.css");

	var Guest = React.createClass({
		getInitialState: function() {
			return {
				processing: false,
				failure: false
			};
		},
		login: function(ev) {
			ev.preventDefault();

			var username = this.refs.username.value;
			var password = this.refs.password.value;

			var xhr = new XMLHttpRequest();
			xhr.onreadystatechange = this.change;
			xhr.onerror = this.error;
			xhr.open("POST", "/system/auth/guest");

			var form = new FormData();
			form.append("user", username);
			form.append("code", password);

			xhr.send(form);
		},
		change: function(ev) {
			console.log(ev);
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
				)));
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

	exports.Login = React.createClass({
		render: function() {
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
					React.DOM.h2(null, "Customer Login:"),
					React.createElement(Guest, {
						url: "/system/auth/guest"
					}),
					React.DOM.h2(null, "Employee Login:"),
					React.createElement(Google)
				))
			);
		}
	});
});
