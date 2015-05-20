//import "/wiki/Wiki.js"
//import "/wiki/Page.js"

(function(Wiki){
	"use strict";

	Wiki.PageProxy = PageProxy;

	function PageProxy(ref, page){
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
	};

	PageProxy.prototype = {

	};
})(Wiki);
