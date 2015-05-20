//import "/wiki/Wiki.js"

(function(Wiki){
	"use strict";

	Wiki.Page = Page;

	function Page(data){
		data = data || {};
		this.owner = data.owner || "";
		this.slug = data.slug || "";
		this.title = data.title || "";
		this.synopsis = data.synopsis || "";
		this.story = data.story || [];
		this.journal = data.journal || [];
	};

	Page.prototype = {

	};
})(Wiki);
