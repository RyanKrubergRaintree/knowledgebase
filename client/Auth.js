package("kb", function (exports) {
	"use strict";

	/*
		Authentication bootup looks like:

		index.html sets up

		window.Provider = {
			community: {type: "*provider.CAS" },
			google: {type: "*provider.Google" },
			guest: {type: "pgdb.GuestLogin" }
		};

		kb.Auth.init extends these providers with login/logout

		when everything has been loaded, it will invoke "loaded" event

		it will tryAutoLogin to different providers, if possible
		if any of them succeeds it will invoke
			"login-success"

		to login with a particular provider you can call:
		window.Provider.guest.login("username", "password");

		if you were already logged in, it will replace the current session,
		however it will not logout from the other session.
	*/

	depends("Session.js");
	depends("util/Notifier.js");

	function Auth(providers, initialSessionInfo) {
		this.notifier_ = new kb.util.Notifier();
		this.notifier_.mixto(this);

		this.loaded = false;
		this.providers = providers;

		this.currentSession = null;
		if (initialSessionInfo) {
			this.currentSession = new kb.Session(initialSessionInfo, function () {});
		}
	}

	Auth.prototype = {
		init: function () {
			var self = this;
			var toload = 1;

			function loaded() {
				toload--;
				if (toload === 0) {
					self.loaded = true;
					self.notifier_.handle({
						type: "loaded"
					});
					if (self.currentSession === null) {
						self.tryAutoLogin();
					}
				}
			}

			for (var name in this.providers) {
				if (!Object.prototype.hasOwnProperty.call(this.providers, name)) {
					continue;
				}
				toload++;
				var data = this.providers[name];
				initializer[data.type](this, name, data, loaded);
			}

			loaded();
		},

		loginSuccess: function (response) {
			if (!response.ok) {
				this.notifier_.handle({
					type: "login-error",
					message: response.text
				});
				this.logoutProviders();
				return;
			}

			var session = new kb.Session(response.json, this.logoutProviders.bind(this));

			this.currentSession = session;
			this.notifier_.handle({
				type: "login-success",
				session: session
			});
		},
		loginError: function (error) {
			this.logoutProviders();
			this.notifier_.handle({
				type: "login-error",
				message: error
			});
		},

		loginTo: function (url, user, code) {
			kb.Session.fetch({
				url: url,
				ondone: this.loginSuccess.bind(this),
				onerror: this.loginError.bind(this),
				body: {
					user: user,
					code: code
				}
			});
		},
		tryAutoLogin: function () {
			for (var name in this.providers) {
				if (!Object.prototype.hasOwnProperty.call(this.providers, name)) {
					continue;
				}
				var provider = this.providers[name];
				if (provider.autologin) {
					provider.autologin();
				}
			}
		},

		logout: function () {
			if (this.currentSession) {
				this.currentSession.logout();
			} else {
				this.logoutProviders();
			}
		},

		// logs out from provider sessions, not from the session
		logoutProviders: function () {
			for (var name in this.providers) {
				if (!Object.prototype.hasOwnProperty.call(this.providers, name)) {
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
		data.login = function (user, code) {
			auth.loginTo("/system/auth/" + name, user, code);
		};

		onloaded();
	}

	function google(auth, name, data, onloaded) {
		if (typeof gapi === "undefined") {
			return;
		}
		data.view = "button";
		data.title = "Google";

		// eslint-disable-next-line no-undef
		gapi.load("auth2", function () {
			if (typeof gapi === "undefined") {
				console.error("Google authentication unavailable.");

				data.login = function () {
					console.error("Google authentication unavailable.");
				};
				data.logout = function () {
					console.error("Google authentication unavailable.");
				};

				data.errored = true;
				onloaded();
				return;
			}

			var hosted_domain = null;
			if (typeof GoogleHostedDomain !== "undefined") {
				// eslint-disable-next-line no-undef
				hosted_domain = GoogleHostedDomain;
			}

			// eslint-disable-next-line no-undef
			var auth2 = gapi.auth2.init({
				hosted_domain: hosted_domain,
				authuser: -1
			});

			var trylogin = function () {
				if (auth2.isSignedIn.get() === true) {
					// check if not logged in
					var user = auth2.currentUser.get();
					var profile = user.getBasicProfile();
					var token = user.getAuthResponse().id_token;

					auth.loginTo("/system/auth/" + name, profile.getEmail(), token);
				}
			};
			auth2.isSignedIn.listen(trylogin);

			data.autologin = function () {
				trylogin();
			};

			data.login = function () {
				if (auth2.isSignedIn.get() === true) {
					trylogin();
				} else {
					auth2.signIn().then(null, auth.loginError.bind(auth));
				}
			};

			data.logout = function () {
				try {
					auth2.signOut();
				} catch (ex) {
					/* empty */
				}
			};

			onloaded();
		});
	}

	function form(auth, name, data, onloaded) {
		data.view = "form";
		data.login = function (user, password) {
			auth.loginTo("/system/auth/" + name, user, password);
		};

		onloaded();
	}

	exports.Auth = new Auth(window.Provider, window.InitialSession);
	exports.Auth.init();
});
