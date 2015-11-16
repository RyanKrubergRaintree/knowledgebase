package("kb.boot", function(exports) {
	"use strict";

	depends("app.css");
	depends("Login.js");

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
				context: {
					user: "Egon Elbre",
					home: "Community=Welcome",
					token: null
				}
			};
		},

		loggedIn: function(user) {

		},
		render: function() {
			var context = this.state.context;
			if (context.token === null) {
				return React.createElement(kb.boot.Login, {
					onLoggedIn: this.loggedIn
				});
			}
			return React.DOM.div({}, "SITE");
		}
	});

	var bootstrap = React.createElement(Bootstrap);
	ReactDOM.render(bootstrap, document.getElementById("boot"));
});
