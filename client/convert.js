import "slug.js";

this.Convert = {};
(function(Convert){
	"use strict";

	// There are several possible links
	// "http://kb.example.com/example"
	// "https://kb.example.com/example"
	// "//kb.example.com/example"
	// "/kb:example" - rooted local URL
	// "kb:Example" - local URL
	Convert.LinkToReference = function(link){
		if(link.indexOf("https://", link) === 0){
			link = link.substr(6);
		} else if(link.indexOf("http://", link) === 0){
			link = link.substr(5);
		}

		link = link.trim();
		// External site:
		// "//kb.example.com/example"
		if((link[0] == "/") && (link[1] == "/") ) {
			return {
				link: Convert.URLToReadable(link),
				url:  link,
				title: Convert.LinkToTitle(link)
			};
		}

		var query = "";
		var q = link.indexOf("?");
		if(q >= 0){
			query = link.substr(q);
			link = link.substr(0, q);
		}

		// remove prefix "/"
		if(link[0] == "/") {
			link = link.substr(1);
		}

		var i = link.lastIndexOf(":")
		var owner = i >= 0 ? link.substr(0,i): "";

		return {
			link: Convert.URLToReadable(link),
			owner: owner,
			url: "/" + Slugify(link) + query,
			title: Convert.LinkToTitle(link),
		};
	}

	Convert.ReferenceToLink = function(ref){
		return Convert.URLToReadable(ref.url);
	};

	Convert.LinkToTitle = function(link){
		var i = Math.max(link.lastIndexOf("/"), link.lastIndexOf(":"));
		link = link.substr(i + 1);
		return link;
	};

	Convert.URLToReadable = function(url){
		var loc = Convert.URLToLocation(url);
		if((typeof loc.origin == "undefined") || (loc.origin == window.location.origin)){
			if(loc.pathname[0] == "/") {
				return loc.pathname + loc.search + loc.hash;
			}
			return "/" + loc.pathname + loc.search + loc.hash;
		}
		return url;
	};

	Convert.URLToLocation = function(url){
		var a = document.createElement("a");
		a.href = url;
		return {
			hash: a.hash,
			search: a.search,
			pathname: a.pathname,
			port: a.port,
			hostname: a.hostname,
			host: a.host,
			password: a.password,
			username: a.username,
			protocol: a.protocol,
			origin: a.origin,
			href: a.href,
		};
	};
})(Convert);
