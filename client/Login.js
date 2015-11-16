package("kb.boot", function(exports) {
	"use strict";

	depends("app.css");
	depends("Login.css");

	exports.Login = React.createClass({
		getInitialState: function() {
			return {
				processing: false,
				failure: false
			};
		},
		loginCustomer: function(ev) {
			ev.preventDefault();
			//
		},
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
					React.DOM.form({
							className: "logins",
							onSubmit: this.loginCustomer
						},
						React.DOM.h2(null, "Customer Login:"),
						React.DOM.table(null, React.DOM.tbody(null,
							React.DOM.tr(null,
								React.DOM.td(null, React.DOM.label({
									htmlFor: "username"
								}, "Username:")),
								React.DOM.td(null, React.DOM.input({
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
									name: "password",
									type: "password",
									tabIndex: 2
								}))
							)
						)),
						React.DOM.div({
								className: "logins"
							},
							React.DOM.h2(null, "Employee Login:"),
							React.DOM.a({
								className: "button",
								href: "#"
							}, "Google")
						)
					)
				))
			);
		}
	});
});


/*
Knowledge Base works best with Google Chrome. Other browsers still have some known issues.

Customer login:

Username:
Login
Password:
Raintree employee login:

Google*/
