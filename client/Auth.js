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

	/**
	 * @param {string} name name of the cookie
	 * @returns {string | null}
	 */
	function getCookie(name) {
		var cookieArr = document.cookie.split(";");

		for (var i = 0; i < cookieArr.length; i++) {
			var cookiePair = cookieArr[i].split("=");
			var cookieKey = cookiePair[0].trim();
			var cookieValue = cookiePair[1];

			if (name == cookieKey) {
				return decodeURIComponent(cookieValue);
			}
		}
		return null;
	}

	function parseJwt(token) {
		var base64Url = token.split(".")[1];
		var base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
		var jsonPayload = decodeURIComponent(
			window
				.atob(base64)
				.split("")
				.map(function (c) {
					return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
				})
				.join("")
		);
		return JSON.parse(jsonPayload);
	}

	var initializer = {
		"*provider.CAS": cas,
		"*provider.Google": initializerGoogle,
		"pgdb.GuestLogin": form
	};

	function cas(auth, name, data, onloaded) {
		data.login = function (user, code) {
			auth.loginTo("/system/auth/" + name, user, code);
		};

		onloaded();
	}

	/**
	 *
	 * @param {Auth} auth
	 * @param {string} name
	 * @param {*} data
	 * @param {() => void} onloaded
	 */
	function initializerGoogle(auth, name, data, onloaded) {
		var gsiInformationElement = document.getElementById("gsi_information");
		if (!gsiInformationElement) {
			console.error("No gsi information given by server");
			return onloaded();
		}
		/** @type {import("./types/globals").GoogleSignInInformation} */
		var gsiInformation = JSON.parse(gsiInformationElement.text);
		if (!gsiInformation.client_id || !gsiInformation.login_uri) {
			console.error("Missing required information for google authentication");
			return onloaded();
		}

		if (
			typeof google === "undefined" ||
			!google.accounts ||
			!google.accounts.id ||
			!google.accounts.id.initialize
		) {
			console.error("Google sign-in library not loaded");
			return onloaded();
		}

		data.view = "google-button";
		data.title = "Google";
		data.cookieName = "gsi_token";

		var autoSelectCookieName = "auto_select";
		var autoSelect = Boolean(JSON.parse(getCookie(autoSelectCookieName)));

		// https://developers.google.com/identity/gsi/web/reference/js-reference#IdConfiguration
		google.accounts.id.initialize({
			client_id: gsiInformation.client_id,
			hd: gsiInformation.hd,
			auto_select: autoSelect,
			login_uri: gsiInformation.login_uri,
			callback: function (response) {
				/** @type {import("./types/globals").ParsedJwtCredentials} */
				var parsedCredential = parseJwt(response.credential);

				auth.loginTo("/system/auth/" + name, parsedCredential.email, response.credential);
				// to speed up consecutive logging
				document.cookie =
					data.cookieName +
					"=" +
					response.credential +
					"; expires=" +
					new Date(parsedCredential.exp * 1000).toUTCString() +
					"; secure";
				document.cookie = autoSelectCookieName + "=true";
			}
		});

		data.autologin = function () {
			var gsiCookie = getCookie(data.cookieName);
			if (!gsiCookie) {
				return google.accounts.id.prompt();
			}

			/** @type {import("./types/globals").ParsedJwtCredentials} */
			var parsedCredential = parseJwt(gsiCookie);
			auth.loginTo("/system/auth/" + name, parsedCredential.email, gsiCookie);
		};

		data.logout = function () {
			var clearCookiePart = "=;expires=Thu, 01 Jan 1970 00:00:00 GMT";
			document.cookie = data.cookieName + clearCookiePart;
			document.cookie = autoSelectCookieName + clearCookiePart;
			// to cancel the one-tap prompt
			google.accounts.id.cancel();
		};
		onloaded();
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
