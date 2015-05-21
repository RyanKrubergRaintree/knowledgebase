//import "/util/Notifier.js"
//import "/wiki/Wiki.js"
//import "/wiki/Page.js"

 (function(Wiki){
	"use strict";

	// Stage represents a staging area where modifications/loading are done.
	Wiki.Stage = Stage;
	function Stage(ref, page){
		this.url = ref.url;
		this.owner = ref.owner;
		this.link = ref.link;
		this.title = ref.title;
		this.key = ref.key;

		page = page || {};
		page.owner = page.owner || ref.owner || "";
		page.slug = page.slug || ref.slug || "";
		page.title = page.title || ref.title || "";

		this.page = new Wiki.Page(page);
		this.notifier = new Notifier();
		this.notifier.mixto(this);
	};

	Stage.prototype = {
		close: function(){
			this.notifier.emit({type: "close"});
		},
	};

})(Wiki);
