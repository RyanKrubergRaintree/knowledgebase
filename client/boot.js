package("kb.boot", function(exports) {
	"use strict";

	depends("app.css");
	depends("Login.js");
	depends("Session.js");

	var Bootstrap = React.createClass({
		componentDidMount: function() {
			var self = this;
			if (typeof Reloader !== "undefined") {
				Reloader.onchange = function() {
					self.forceUpdate();
				};
			}
		},
		componentWillUnmount: function() {},

		getInitialState: function() {
			return {
				session: null
			};
		},

		loggedIn: function(session) {
			this.setState({
				session: new kb.Session(session)
			});
		},
		render: function() {
			var session = this.state.session;
			if (session === null) {
				return React.createElement(kb.boot.Login, {
					onSuccess: this.loggedIn
				});
			}

			return React.DOM.div({},
				"Logged in as: ",
				JSON.stringify(session)
			);
		}
	});

	var bootstrap = React.createElement(Bootstrap);
	ReactDOM.render(bootstrap, document.getElementById("boot"));
});
