//import "/util/ObjectId.js"
//import "/util/Notifier.js"
//import "/kb/KB.js"
//import "/kb/Page.js"

KB.Stage = (function(){
	"use strict";

	// Stage represents a staging area where modifications/loading are done.
	function Stage(ref, page){
		this.id = GenerateID();

		this.url = ref.url;
		this.link = ref.link;
		this.title = ref.title;

		page = page || {};
		page.title = page.title || ref.title || "";

		this.page = new KB.Page(page);
		this.notifier = new Notifier();
		this.notifier.mixto(this);

		this.state_ = "";

		var self = this;
		window.setTimeout(function(){
			self.requestPage();
		});

		this.lastStatus = 200;
		this.lastError = "";
	};

	Stage.prototype = {
		set state(value){
			this.state_ = value;
			this.changed();
		},

		get state(){
			return this.state_;
		},

		close: function(){
			this.state = "closed";
			this.notifier.handle({type: "closed", stage: this});
		},
		changed: function(){
			this.notifier.emit({type: "changed", stage: this});
		},

		requestPage: function(){
			this.state = "requesting";

			var xhr = new XMLHttpRequest();
			xhr.withCredentials = true;
			var self = this;

			xhr.addEventListener('load', function(ev){
				var data = JSON.parse(xhr.response),
					page = new KB.Page(data);
				self.page = page;
				self.state = "loaded";
			}, false);

			xhr.addEventListener('error', function(ev){
				console.error(ev);
			}, false);

			xhr.open('GET', this.url, true);
			xhr.setRequestHeader('Accept', 'application/json');
			xhr.send();
		}
	};

	return Stage;
})();
