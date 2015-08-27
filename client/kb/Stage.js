package('kb', function(exports){
	'use strict';

	depends('util/Notifier.js');
	depends('Page.js');
	depends('Tracking.js');

	depends('util/ParseJSON.js');

	function success(xhr){
		return (200 <= xhr.status) && (xhr.status < 300);
	}

	function bindready(xhr, fn, self){
		return function(){
			if(xhr.readyState !== 4){ return; }
			fn.call(self, xhr);
		}
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
		}
	};

	// Stage represents a staging area where modifications/loading are done.
	exports.Stage = Stage;
	function Stage(ref, page){
		this.id = GenerateID();

		this.creating = (ref.url === null) || (ref.url === '');
		this.url = ref.url;
		this.link = ref.link;
		this.title = ref.title;
		this.allowed = ['GET', 'HEAD'];

		page = page || {};
		page.title = page.title || ref.title || '';

		this.page = new kb.Page(page);
		this.editing = new Editing(this);

		this.notifier = new kb.util.Notifier();
		this.notifier.mixto(this);

		this.state = '';

		this.lastStatus = 200;
		this.lastStatusText = '';
		this.lastError = '';

		this.patching_ = false;
		this.patches_ = [];

		this.wide = false;
	}

	Stage.prototype = {
		close: function(){
			this.state = 'closed';
			this.changed();
			this.notifier.handle({type: 'closed', stage: this});
		},
		changed: function(){
			this.notifier.emit({type: 'changed', stage: this});
		},
		urlChanged: function(){
			this.notifier.emit({type: 'urlChanged', stage: this});
		},

		wideChanged: function(){
			this.notifier.emit({type: 'widthChanged', stage: this});
		},
		expand: function(){
			this.wide = true;
			this.wideChanged();
		},
		collapse: function(){
			this.wide = false;
			this.wideChanged();
		},

		canCreate: function(){
			return this.allowed.indexOf('PUT') >= 0;
		},
		canModify: function(){
			return this.allowed.indexOf('PATCH') >= 0;
		},
		canDestroy: function(){
			return this.allowed.indexOf('DELETE') >= 0;
		},

		updateStatus_: function(xhr){
			var allowed = xhr.getResponseHeader('Allow');
			if(typeof allowed === 'string'){
				this.allowed = allowed.split(',').map(function(v){ return v.trim(); });
			}

			var ok = success(xhr);
			this.state = 'loaded';
			if(!ok){
				this.state = 'error';
				if(xhr.status === 404){
					this.state = 'not-found';
					if(this.canCreate()){
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

			// var version = this.page.version;
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
				xhr.onreadystatechange = bindready(xhr, this.patchDone_, this);
				xhr.onerror = this.patchError_.bind(this);

				xhr.open('PATCH', this.url, true);

				xhr.setRequestHeader('Accept', 'application/json');
				xhr.setRequestHeader('Content-Type', 'application/json');

				xhr.send(JSON.stringify(patch));
			}
		},
		patchDone_: function(ev){
			var xhr = ev.target;
			if(xhr.readyState !== 4){ return; }

			this.patching_ = false;

			if(!this.updateStatus_(xhr)){
				//TODO: don't drop changes in case of errors
				this.patches_ = [];
				this.patching_ = false;
				this.pull();
				return;
			}
			this.nextPatch_();
		},
		patchError_: function(/* ev */){
			this.patches_ = [];
			this.patching_ = false;
			this.pull();
		},

		refresh: function(){
			this.pull();
		},
		pull: function(){
			if(this.url == null){ return; }

			this.state = 'loading';
			this.changed();

			var xhr = new XMLHttpRequest();
			xhr.onreadystatechange = bindready(xhr, this.pullDone_, this);
			xhr.onerror = this.pullError_.bind(this);

			xhr.open('GET', this.url, true);
			xhr.setRequestHeader('Accept', 'application/json');
			xhr.send();
		},
		pullDone_: function(xhr){
			if(!this.updateStatus_(xhr)){
				this.changed();
				return;
			}

			var data = kb.util.ParseJSON(xhr.responseText),
			page = new kb.Page(data);
			if(xhr.responseURL){
				if(this.url !== xhr.responseURL){
					this.url = xhr.responseURL;
					this.urlChanged();
				}
			}


			this.page = page;
			this.state = 'loaded';
			this.changed();

			kb.TrackPageView(this.url, this.page.title);
		},
		pullError_: function(/* ev */){
			this.state = 'failed';
			this.lastStatus = 'failed';
			this.lastStatusText = '';
			this.lastError = '';
			this.changed();
		},

		create: function(){
			if(!this.creating){ return; }
			this.url = '/' + this.link;
			this.urlChanged();

			var xhr = new XMLHttpRequest();
			xhr.onreadystatechange = bindready(xhr, this.createDone_, this);
			xhr.onerror = this.createError_.bind(this);

			xhr.open('PUT', this.url, true);
			xhr.setRequestHeader('Accept', 'application/json');
			xhr.setRequestHeader('Content-Type', 'application/json');
			xhr.send(JSON.stringify({
				title: this.title,
				slug: this.link,
				story: [
					{ id: GenerateID(), type: 'tags' },
					{ id: GenerateID(), type: 'factory'}
				]
			}));

			this.changed();
		},
		createDone_: function(xhr){
			if(!this.updateStatus_(xhr)){
				this.changed();
				return;
			}
			this.creating = false;
			this.state = 'created';
			this.refresh();
		},
		createError_: function(/* ev */){
			this.state = 'failed';
			this.lastStatus = 'failed';
			this.lastStatusText = '';
			this.lastError = '';
			this.changed();
		},

		destroy: function(){
			if(this.url == null){ return; }

			var xhr = new XMLHttpRequest();
			xhr.open('DELETE', this.url, true);

			xhr.onreadystatechange = bindready(xhr, this.destroyDone_, this);
			xhr.onerror = this.destroyError_.bind(this);

			xhr.send();
			return xhr;
		},
		destroyDone_: function(ev){
			this.updateStatus_(ev.target);
			this.changed();
			this.pull();
		},
		destroyError_: function(/* ev */){
			this.state = 'failed';
			this.lastStatus = 'failed';
			this.lastStatusText = '';
			this.lastError = '';
			this.changed();
		}
	};
});