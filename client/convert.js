import "slug.js";

this.Convert = {};
(function(Convert){
	"use strict";

	function trimProtocol(link){
		if(link.indexOf("https://", link) === 0){
			return link.substr(6).trim();
		} else if(link.indexOf("http://", link) === 0){
			return link.substr(5).trim();
		}
		return link.trim();
	}

	function trimLeadingSlashes(link){
		// remove prefix "/"
		while(link[0] == "/") {
			link = link.substr(1);
		}
		return link;
	}

	// There are several possible links
	// "http://kb.example.com/example"
	// "https://kb.example.com/example"
	// "//kb.example.com/example"
	// "/kb:example" - rooted local URL
	// "kb:Example" - local URL
	Convert.LinkToReference = function(link){
		link = trimProtocol(link);
		// External site:
		// "//kb.example.com/example"
		if((link[0] == "/") && (link[1] == "/") ) {
			return {
				link: Convert.URLToReadable(link),
				url:  link,
				title: Convert.LinkToTitle(link)
			};
		}

		link = trimLeadingSlashes(link);
		var i = link.indexOf(":")
		var owner = i >= 0 ? link.substr(0,i): "";

		return {
			link: Convert.URLToReadable(link),
			owner: owner,
			url: "/" + Slugify(link),
			title: Convert.LinkToTitle(link),
		};
	}

	Convert.ReferenceToLink = function(ref){
		return Convert.URLToReadable(ref.url);
	};

	Convert.LinkToTitle = function(link){
		link = trimProtocol(link);
		link = trimLeadingSlashes(link);
		var i = Math.max(link.lastIndexOf("/"), link.indexOf(":"));
		link = link.substr(i + 1);
		return link;
	};

	Convert.LinkToOwner = function(link){
		link = Convert.URLToReadable(link);
		link = trimProtocol(link);
		link = trimLeadingSlashes(link);
		var i = link.indexOf(":");
		link = link.substr(0, i);
		return link.trim().toLowerCase();
	}

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
