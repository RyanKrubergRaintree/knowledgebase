package("kb", function (exports) {
	"use strict";

	depends("util/Notifier.js");
	depends("Page.js");
	depends("Tracking.js");

	depends("util/ParseJSON.js");

	function Editing(stage) {
		this.edited = false;
		this.stage = stage;
		this.items = {};
	}
	Editing.prototype = {
		start: function (id) {
			this.edited = true;
			this.items[id] = true;
			this.pack();
			this.stage.changed();
		},
		stop: function (id) {
			delete this.items[id];
			this.pack();
			this.stage.changed();
		},
		// removes all non-existing items
		pack: function () {
			var t = {};
			var page = this.stage.page;
			for (var i = 0; i < page.story.length; i++) {
				var item = page.story[i];
				if (this.items[item.id]) {
					t[item.id] = true;
				}
			}
			this.items = t;
		},
		clear: function () {
			this.items = {};
			this.changed();
		},
		item: function (id) {
			return this.items[id];
		}
	};

	// Stage represents a staging area where modifications/loading are done.
	exports.Stage = Stage;

	function Stage(session, ref, page) {
		this.session_ = session;
		this.id = GenerateID();

		this.creating = ref.url === null || ref.url === "";
		this.link = ref.link;
		this.title = ref.title;
		this.allowed = ["GET", "HEAD"];

		if (ref.url !== null) {
			this.url = ref.url.replace("/#", "");
		} else {
			this.url = ref.url;
		}

		page = page || {};
		page.title = page.title || ref.title || "";

		this.customClassName = kb.convert.URLGetQueryParam(ref.url, "className");

		this.page = new kb.Page(page);
		this.editing = new Editing(this);

		this.notifier = new kb.util.Notifier();
		this.notifier.mixto(this);

		this.state = "";

		this.lastStatus = 200;
		this.lastStatusText = "";
		this.lastError = "";

		this.patching_ = false;
		this.patches_ = [];

		this.wide = true;
	}

	Stage.prototype = {
		close: function () {
			this.changed();
			this.notifier.handle({
				type: "closed",
				stage: this
			});
		},
		changed: function (loaded) {
			this.notifier.emit({
				type: "changed",
				stage: this,
				loaded: loaded === true
			});
		},
		urlChanged: function () {
			this.notifier.emit({
				type: "urlChanged",
				stage: this
			});
		},

		wideChanged: function () {
			this.notifier.emit({
				type: "widthChanged",
				stage: this
			});
		},
		expand: function () {
			this.wide = true;
			this.wideChanged();
		},
		collapse: function () {
			this.wide = false;
			this.wideChanged();
		},

		canCreate: function () {
			return this.allowed.indexOf("PUT") >= 0;
		},
		canModify: function () {
			return this.allowed.indexOf("POST") >= 0;
		},
		canDestroy: function () {
			return this.allowed.indexOf("DELETE") >= 0;
		},
		canViewHistory: function () {
			return this.allowed.indexOf("OVERWRITE") >= 0;
		},

		updateStatus_: function (response) {
			var allowed = response.xhr.getResponseHeader("Allow");
			if (typeof allowed === "string") {
				this.allowed = allowed.split(",").map(function (v) {
					return v.trim();
				});
			}

			this.state = "loaded";
			if (!response.ok) {
				this.state = "error";
				if (response.xhr.status === 404) {
					this.state = "not-found";
					if (this.canCreate()) {
						this.creating = true;
					}
				}
				if (response.xhr.status === 204) {
					this.close();
				}
			}

			this.lastStatus = response.status;
			this.lastStatusText = response.statusText;
			this.lastError = response.text;

			return response.ok;
		},

		patch: function (op) {
			if (this.url === null || this.url === "") {
				return;
			}

			// var version = this.page.version;
			this.page.apply(op);

			this.patches_.push(op);
			this.nextPatch_();

			this.changed(true);
		},
		nextPatch_: function () {
			if (this.patching_) {
				return;
			}
			var patch = this.patches_.shift();
			if (patch) {
				this.patching_ = true;
				this.session_.fetch({
					method: "POST",
					url: this.url,
					ondone: this.patchDone_.bind(this),
					onerror: this.patchError_.bind(this),
					headers: {
						Accept: "application/json",
						"Content-Type": "application/json"
					},
					body: JSON.stringify(patch)
				});
			}
		},
		patchDone_: function (response) {
			this.patching_ = false;

			if (!this.updateStatus_(response)) {
				//TODO: don't drop changes in case of errors
				this.patches_ = [];
				this.patching_ = false;
				this.pull();
				return;
			}
			this.nextPatch_();
		},
		patchError_: function (/* ev */) {
			this.patches_ = [];
			this.patching_ = false;
			this.pull();
		},

		refresh: function () {
			this.pull();
		},
		pull: function () {
			if (this.url === null || this.url === "") {
				return;
			}

			this.state = "loading";
			this.changed(false);

			this.session_.fetch({
				method: "GET",
				url: this.url,
				ondone: this.pullDone_.bind(this),
				onerror: this.pullError_.bind(this),
				headers: {
					Accept: "application/json"
				}
			});
		},
		pullDone_: function (response) {
			if (!this.updateStatus_(response)) {
				this.changed(true);
				return;
			}

			var data = kb.util.ParseJSON(response.text),
				page = new kb.Page(data);

			if (this.url !== response.url) {
				this.url = response.url;
				this.urlChanged();
			}

			this.page = page;
			this.state = "loaded";
			this.changed(true);

			var niceurl = kb.convert.URLToReadable(this.url);
			kb.TrackPageView(niceurl, this.page.title);
		},
		pullError_: function (/* ev */) {
			this.state = "failed";
			this.lastStatus = "failed";
			this.lastStatusText = "";
			this.lastError = "";
			this.changed(false);
		},

		create: function () {
			if (!this.creating) {
				return;
			}
			this.url = "/" + this.link;
			this.urlChanged();

			this.session_.fetch({
				method: "PUT",
				url: this.url,
				ondone: this.createDone_.bind(this),
				onerror: this.createError_.bind(this),
				headers: {
					Accept: "application/json",
					"Content-Type": "application/json"
				},
				body: JSON.stringify({
					title: this.title,
					slug: this.link,
					story: [
						{
							id: GenerateID(),
							type: "tags"
						},
						{
							id: GenerateID(),
							type: "factory"
						}
					]
				})
			});

			this.changed(false);
		},
		createDone_: function (xhr) {
			if (!this.updateStatus_(xhr)) {
				this.changed();
				return;
			}
			this.creating = false;
			this.state = "created";
			this.refresh();
		},
		createError_: function (/* ev */) {
			this.state = "failed";
			this.lastStatus = "failed";
			this.lastStatusText = "";
			this.lastError = "";
			this.changed(false);
		},

		destroy: function () {
			if (this.url === null || this.url === "") {
				return;
			}

			this.session_.fetch({
				method: "DELETE",
				url: this.url,
				ondone: this.destroyDone_.bind(this),
				onerror: this.destroyError_.bind(this)
			});
		},
		destroyDone_: function (response) {
			this.updateStatus_(response);
			this.changed(false);
			this.pull();
		},
		destroyError_: function (/* ev */) {
			this.state = "failed";
			this.lastStatus = "failed";
			this.lastStatusText = "";
			this.lastError = "";
			this.changed(false);
		}
	};
});
