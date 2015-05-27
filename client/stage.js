import "util/ObjectId.js";
import "util/notifier.js";
import "kb.js";
import "page.js";

KB.Stage = (function(){
	"use strict";

	function success(xhr){
		return (200 <= xhr.status) && (xhr.status < 300);
	}

	function Editing(stage){
		this.stage = stage;
		this.items = {};
	}
	Editing.prototype = {
		start: function(id){
			this.items[id] = true;
			this.pack();
			this.stage.changed();
		},
		stop: function(id){
			delete this.items[id];
			this.pack();
			this.stage.changed();
		},
		// removes all non-existing items
		pack: function(){
			var t = {};
			var page = this.stage.page;
			for(var i = 0; i < page.story.length; i++){
				var item = page.story[i];
				if(this.items[item.id]){
					t[item.id] = true;
				}
			}
			this.items = t;
		},
		clear: function(){
			this.items = {};
			this.changed();
		},
		item: function(id){
			return this.items[id];
		},
	};

	// Stage represents a staging area where modifications/loading are done.
	function Stage(ref, page){
		this.id = GenerateID();

		this.creating = ref.url == null;
		this.url = ref.url;
		this.link = ref.link;
		this.title = ref.title;

		page = page || {};
		page.title = page.title || ref.title || "";

		this.page = new KB.Page(page);
		this.editing = new Editing(this);

		this.notifier = new Notifier();
		this.notifier.mixto(this);

		this.state = "";

		this.lastStatus = 200;
		this.lastStatusText = "";
		this.lastError = "";

		this.patching_ = false;
		this.patches_ = [];
	};

	Stage.prototype = {
		close: function(){
			this.state = "closed";
			this.changed();
			this.notifier.handle({type: "closed", stage: this});
		},
		changed: function(){
			this.notifier.emit({type: "changed", stage: this});
		},
		urlChanged: function(){
			this.notifier.emit({type: "urlChanged", stage: this});
		},

		get canCreate(){
			//TODO: implement
			return true;
		},
		get canModify(){
			//TODO: implement
			return true;
		},

		updateStatus_: function(xhr){
			var ok = success(xhr);
			this.state = 'loaded';
			if(!ok){
				this.state = 'error';
				if(xhr.status === 404){
					this.state = 'not-found';
					if(this.canCreate){
						this.creating = true;
					}
				}
			}

			this.lastStatus = xhr.status;
			this.lastStatusText = xhr.statusText;
			this.lastError = xhr.responseText;

			return ok;
		},

		patch: function(op){
			if(this.url == null){ return; }

			var version = this.page.version;
			this.page.apply(op);

			this.patches_.push(op);
			this.nextPatch_();

			this.changed();
		},
		nextPatch_: function(){
			if(this.patching_){ return; }
			var patch = this.patches_.shift();
			if(patch){
				this.patching_ = true;

				var xhr = new XMLHttpRequest();
				xhr.withCredentials = true;
				xhr.addEventListener('load', this.patchDone_.bind(this), false);
				xhr.addEventListener('error', this.patchError_.bind(this), false);

				xhr.open("PATCH", this.url, true);

				xhr.setRequestHeader('Accept', 'application/json');
				xhr.setRequestHeader('Content-Type', 'application/json');

				xhr.send(JSON.stringify(patch));
			}
		},
		patchDone_: function(ev){
			this.patching_ = false;
			var xhr = ev.target;
			if(!this.updateStatus_(xhr)){
				//TODO: don't drop changes in case of errors
				this.patches_ = [];
				this.patching_ = false;
				this.pull();
				return;
			}
			this.nextPatch_();
		},
		patchError_: function(ev){
			this.patches_ = [];
			this.patching_ = false;
			this.pull();
		},


		pull: function(){
			if(this.url == null){ return; }

			this.state = "loading";
			this.changed();

			var xhr = new XMLHttpRequest();
			xhr.withCredentials = true;
			xhr.addEventListener('load', this.pullDone_.bind(this), false);
			xhr.addEventListener('error', this.pullError_.bind(this), false);

			xhr.open('GET', this.url, true);
			xhr.setRequestHeader('Accept', 'application/json');
			xhr.send();
		},
		pullDone_: function(ev){
			var xhr = ev.target;
			if(!this.updateStatus_(xhr)){
				this.changed();
				return;
			}

			var data = JSON.parse(xhr.response),
			page = new KB.Page(data);
			if(xhr.responseURL){
				if(this.url != xhr.responseURL){
					this.url = xhr.responseURL;
					this.urlChanged();
				}
			}
			this.page = page;
			this.state = "loaded";
			this.changed();
		},
		pullError_: function(ev){
			this.state = "failed";
			this.lastStatus = "failed";
			this.lastStatusText = "";
			this.lastError = "";
			this.changed();
		},

		create: function(){
			if(!this.creating){ return; }
			this.url = "/" + this.link;
			this.urlChanged();

			var xhr = new XMLHttpRequest();
			xhr.withCredentials = true;
			xhr.addEventListener('load', this.createDone_.bind(this), false);
			xhr.addEventListener('error', this.createError_.bind(this), false);

			xhr.open('PUT', this.url, true);
			xhr.setRequestHeader('Accept', 'application/json');
			xhr.setRequestHeader('Content-Type', 'application/json');
			xhr.send(JSON.stringify({
				title: this.title,
				slug: this.link
			}));

			this.changed();
		},
		createDone_: function(ev){
			var xhr = ev.target;
			if(!this.updateStatus_(xhr)){
				this.changed();
				return;
			}
			this.creating = false;

			var data = JSON.parse(xhr.response),
			page = new KB.Page(data);
			if(xhr.responseURL){
				this.url = xhr.responseURL;
			}
			this.page = page;
			this.state = "loaded";
			this.changed();
		},
		createError_: function(ev){
			this.state = "failed";
			this.lastStatus = "failed";
			this.lastStatusText = "";
			this.lastError = "";
			this.changed();
		},

		destroy: function(){
			if(this.url == null){ return; }

			var xhr = new XMLHttpRequest();
			xhr.withCredentials = true;
			xhr.open('DELETE', jsonurl(this.url), true);

			xhr.addEventListener('load', this.destroyDone_.bind(this), false);
			xhr.addEventListener('error', this.destroyError_.bind(this), false);

			xhr.send();
			return xhr;
		},
		destroyDone_: function(ev){
			this.updateStatus_(ev.target);
			this.changed();
			this.pull();
		},
		destroyError_: function(ev){
			this.state = "failed";
			this.lastStatus = "failed";
			this.lastStatusText = "";
			this.lastError = "";
			this.changed();
		}

	};

	return Stage;
})();