package("kb", function(exports) {
	"use strict";

	depends("Session.js");
	depends("Notifier.js");

	function Auth(providers) {
		this.notifier_ = new kb.util.Notifier();
		this.notifier_.mixto(this);

		this.loaded = false;
		this.providers = providers;
	}

	Auth.prototype = {
		init: function() {
			var self = this;
			var toload = 1;

			function loaded() {
				toload--;
				if (toload === 0) {
					self.loaded = true;
					self.notifier_.handle({
						type: "loaded"
					});
					self.tryAutoLogin();
				}
			}

			for (var name in this.providers) {
				if (!this.providers.hasOwnProperty(name)) {
					continue;
				}
				var data = this.providers[name];
				initializer[data.type](this, name, data, loaded);
			}
		},

		loginSuccess: function(response) {
			if (!response.ok) {
				this.notifier_.handle({
					type: "login-error",
					message: response.text
				});
				this.logout(true);
				return;
			}

			var session = new kb.Session(
				response.json,
				this.logout.bind(this)
			);

			this.notifier_.handle({
				type: "login-success",
				session: session
			});
		},
		loginError: function(error) {
			this.logout();
			this.notifier_.handle({
				type: "login-error",
				message: error
			});
		},

		loginTo: function(url, user, code) {
			var form = new FormData();
			form.append("user", user);
			form.append("code", code);

			kb.Session.fetch({
				url: url,
				ondone: this.loginSuccess.bind(this),
				onerror: this.loginError.bind(this),
				body: form
			});
		},
		tryAutoLogin: function() {

		},

		// logs out from provider sessions, not from the session
		logout: function() {
			for (var name in this.providers) {
				if (!this.providers.hasOwnProperty(name)) {
					continue;
				}
				var provider = this.providers[name];
				if (provider.logout) {
					provider.logout();
				}
			}
		}
	};

	var initializer = {
		"*provider.CAS": cas,
		"*provider.Google": google,
		"pgdb.GuestLogin": form
	};

	function cas(auth, name, data, onloaded) {
		data.login = function(user, code) {
			auth.loginTo("/system/auth/" + name, user, code);
		};
		onloaded();
	}

	function google(auth, name, data, onloaded) {
		data.view = "button";
		data.title = "Google";

		gapi.load("auth2", function() {
			data.autologin = function() {

			};

			data.login = function() {

			};

			data.logout = function() {
				if (gapi.auth2) {
					var auth = gapi.auth2.getAuthInstance();
					if (auth) {
						try {
							auth.signOut();
						} catch (ex) {}
					}
				}
			};
			onloaded();
		});
	}

	function form(auth, name, data, onloaded) {
		data.view = "form";
		data.login = function(user, password) {
			auth.loginTo("/system/auth/" + name, user, password);
		};
		onloaded();
	}

	exports.Auth = new Auth(window.Provider);
	exports.Auth.init();
});
