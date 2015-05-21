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
		this.owner = ref.owner;
		this.link = ref.link;
		this.title = ref.title;

		page = page || {};
		page.owner = page.owner || ref.owner || "";
		page.slug = page.slug || ref.slug || "";
		page.title = page.title || ref.title || "";

		this.page = new KB.Page(page);
		this.notifier = new Notifier();
		this.notifier.mixto(this);
	};

	Stage.prototype = {
		close: function(){
			this.notifier.emit({type: "close"});
		},

		resolveLinks: function(text){
			return text;
		}
	};

	return Stage;
})();
