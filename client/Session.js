package("kb", function(exports) {
	"use strict";

	exports.Session = Session;

	function Session(context) {
		context = context || {};
		this.user = context.user || {
			id: "",
			email: "",
			name: "",
			company: "",
			admin: false
		};
		this.home = "Community=Welcome";
		this.branch = context.branch || "10.2.600";
		this.token = context.token || null;
	}
	Session.fetch = function(opts) {
		(new Session()).fetch(opts);
	};

	Session.prototype = {
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

			var xhr = new XMLHttpRequest();
			xhr.onreadystatechange = function() {
				if (xhr.readyState !== 4) {
					return;
				}

				//TODO: add authentication failure error

				opts.ondone({
					get json() {
						return JSON.parse(xhr.responseText);
					},
					url: xhr.responseURL || opts.url,
					status: xhr.status,
					ok: xhr.status === 200,
					statusText: xhr.statusText,
					text: xhr.responseText,
					xhr: xhr
				});
			};

			xhr.onerror = function(err) {
				opts.onerror(err);
			};

			xhr.open(opts.method, opts.url);

			for (var name in opts.headers) {
				if (!opts.headers.hasOwnProperty(name)) {
					continue;
				}
				xhr.setRequestHeader(name, opts.headers[name]);
			}

			if ((typeof opts.body === "undefined") || (opts.body === null)) {
				xhr.send();
			} else {
				xhr.send(opts.body);
			}

			return xhr;
		}
	};
});
