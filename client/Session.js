package("kb", function(exports) {
	"use strict";

	depends("util/Notifier.js");

	exports.Session = Session;

	function Session(context, logoutProvider) {
		this.notifier_ = new kb.util.Notifier();
		this.notifier_.mixto(this);

		context = context || {};
		this.user = context.user || {
			id: "",
			email: "",
			name: "",
			company: "",
			admin: false
		};

		var params = context.params || {};
		this.pages = context.pages || [];
		this.home = params.home || "Community=Welcome";
		this.branch = params.branch || "";
		this.token = context.token || null;
		this.filter = params.branch || "10.2.700";

		this.logoutProvider_ = logoutProvider;
	}
	Session.fetch = function(opts) {
		(new Session()).fetch(opts);
	};

	Session.prototype = {
		logout: function() {
			this.fetch({
				url: "/system/auth/logout"
			});

			this.logoutProvider_();
			this.notifier_.emit({
				type: "session-finished",
				error: ""
			});
		},
		fetch: function(opts) {
			if (typeof opts.url === "undefined") {
				throw new Error("No url defined.");
			}

			opts.method = opts.method || "POST";
			opts.ondone = opts.ondone || function() {};
			opts.onerror = opts.onerror || function() {};

			opts.headers = opts.headers || {};
			if (this.token) {
				opts.headers["X-Auth-Token"] = opts.headers["X-Auth-Token"] || this.token;
			}

			if (["GET", "PUT", "POST", "DELETE"].indexOf(opts.method) < 0) {
				throw new Error("Invalid method: " + opts.method);
			}

			var self = this;

			var xhr = new XMLHttpRequest();
			xhr.onreadystatechange = function() {
				if (xhr.readyState !== 4) {
					return;
				}

				var json = null;
				try {
					json = JSON.parse(xhr.responseText);
				} catch (err) {
					json = null;
				}

				var response = {
					json: json,
					url: xhr.responseURL || opts.url,
					status: xhr.status,
					ok: xhr.status === 200,
					statusText: xhr.statusText,
					text: xhr.responseText,
					xhr: xhr
				};

				opts.ondone(response);

				if (response.status === 401) {
					self.notifier_.emit({
						type: "session-finished",
						error: response.text
					});
					return;
				}
			};

			xhr.onerror = function(err) {
				opts.onerror(err);
			};

			xhr.open(opts.method, opts.url);

			if (this.filter) {
				xhr.setRequestHeader("X-Filter", this.filter);
			}

			for (var name in opts.headers) {
				if (!opts.headers.hasOwnProperty(name)) {
					continue;
				}
				xhr.setRequestHeader(name, opts.headers[name]);
			}

			if ((typeof opts.body === "undefined") || (opts.body === null)) {
				xhr.send();
			} else if (typeof opts.body === "string") {
				xhr.send(opts.body);
			} else {
				var pairs = [];
				for (name in opts.body) {
					if (opts.body.hasOwnProperty(name)) {
						pairs.push(
							encodeURIComponent(name) + "=" +
							encodeURIComponent(opts.body[name]));
					}
				}

				var encoded = pairs.join("&").replace(/%20/g, "+");
				xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");

				xhr.send(encoded);
			}

			return xhr;
		}
	};
});
