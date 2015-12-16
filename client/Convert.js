package("kb.convert", function(exports) {
	"use strict";

	depends("unicode/identifier.js");
	depends("unicode/runename.js");

	function trimProtocol(link) {
		if (link.indexOf("https://", link) === 0) {
			return link.substr(6).trim();
		} else if (link.indexOf("http://", link) === 0) {
			return link.substr(5).trim();
		}
		return link.trim();
	}

	TestCase("trimProtocol", function(assert) {
		assert.equal(trimProtocol(""), "");
		assert.equal(trimProtocol("http://xyz"), "//xyz");
		assert.equal(trimProtocol("https://xyz/zwy"), "//xyz/zwy");
	});

	function trimLeadingSlashes(link) {
		// remove prefix "/"
		while (link.charAt(0) === "/") {
			link = link.substr(1);
		}
		return link;
	}

	TestCase("trimLeadingSlashes", function(assert) {
		assert.equal(trimLeadingSlashes(""), "");
		assert.equal(trimLeadingSlashes("xyz"), "xyz");
		assert.equal(trimLeadingSlashes("///xyz/zwy"), "xyz/zwy");
	});


	// TextToSlug converts text to a slug
	//
	// * numbers and "/" are left intact
	// * letters will be lowercased (if possible)
	// * "-", ",", ".", " ", "_" will be converted to "-"
	// * other symbols or punctuations will be converted to html entity reference name
	//   (if there exists such reference name)
	// * everything else will be converted to "-"
	//
	// Example:
	//   "&Hello_世界/+!" ==> "amp-hello-世界/plus-excl"
	//   "Hello  World  /  Test" ==> "hello-world/test"
	exports.TextToSlug = TextToSlug;

	function TextToSlug(title) {
		var cutdash = true,
			emitdash = false;

		var slug = "";

		for (var i = 0; i < title.length; i += 1) {
			var r = title.charAt(i);
			if (kb.unicode.IsIdent(r)) {
				if (emitdash && !cutdash) {
					slug += "-";
				}
				slug += r.toLowerCase();

				emitdash = false;
				cutdash = false;
				continue;
			}
			if ((r === "/") || (r === "=")) {
				if ((slug.length == 0) || (slug[slug.length - 1] != r)) {
					slug += r;
				}
				emitdash = false;
				cutdash = true;
			} else if ((r === "-") || (r === ",") || (r === ".") || (r === " ") || (r === "_")) {
				emitdash = true;
			} else {
				var name = kb.unicode.RuneName[r];
				if (name) {
					if (!cutdash) {
						slug += "-";
					}
					slug += name;
					cutdash = false;
				}
				emitdash = true;
			}
		}

		if (slug.length === 0) {
			return "-";
		}

		return slug;
	}

	TestCase("TextToSlug", function(assert) {
		assert.equal(TextToSlug(""), "-");
		assert.equal(TextToSlug("&Hello_世界/+!"), "amp-hello-世界/plus-excl");
		assert.equal(TextToSlug("Hello  World  ////  Test"), "hello-world/test");
		assert.equal(TextToSlug("alpha====beta"), "alpha=beta");
	});

	// There are several possible links
	// "http://kb.example.com/example"
	// "https://kb.example.com/example"
	// "//kb.example.com/example"
	// "/kb=example" - rooted local URL
	// "kb=Example" - local URL
	exports.LinkToReference = LinkToReference;

	function LinkToReference(link) {
		link = trimProtocol(link);
		// External site:
		// "//kb.example.com/example"
		if ((link[0] === "/") && (link[1] === "/")) {
			return {
				link: URLToReadable(link),
				url: link,
				title: LinkToTitle(link)
			};
		}

		link = trimLeadingSlashes(link);
		var i = link.indexOf("=");
		var owner = i >= 0 ? link.substr(0, i) : "";

		return {
			link: URLToReadable(link),
			owner: owner,
			url: URLToReadable(link),
			title: LinkToTitle(link)
		};
	}

	exports.ReferenceToLink = ReferenceToLink;

	function ReferenceToLink(ref) {
		return URLToReadable(ref.url);
	}

	exports.LinkToTitle = LinkToTitle;

	function LinkToTitle(link) {
		link = trimProtocol(link);
		link = trimLeadingSlashes(link);
		var i = Math.max(link.lastIndexOf("/"), link.indexOf("="));
		link = link.substr(i + 1);
		return link;
	}

	exports.LinkToOwner = LinkToOwner;

	function LinkToOwner(link) {
		link = URLToReadable(link);
		link = trimProtocol(link);
		link = trimLeadingSlashes(link);

		var i = link.indexOf("=");
		link = link.substr(0, i);
		return link.trim().toLowerCase();
	}

	exports.URLToReadable = URLToReadable;

	function URLToReadable(url) {
		var loc = URLToLocation(url);
		if ((loc.host === "") || (loc.host === window.location.host)) {
			return loc.path + loc.query + loc.fragment;
		} else {
			return "//" + loc.host + loc.path + loc.query + loc.fragment;
		}
	}
	TestCase("URLToReadable", function(assert) {
		assert.equal(URLToReadable(""), "/");
		assert.equal(URLToReadable("/hello-world"), "/hello-world");
		assert.equal(URLToReadable("/hello-world/test#"), "/hello-world/test");
		assert.equal(URLToReadable("https://" + window.location.host + "/hello-world"), "/hello-world");
		assert.equal(URLToReadable("http://unknown.com/hello-world"), "//unknown.com/hello-world");
		assert.equal(URLToReadable("http://unknown.com:241/hello-world"), "//unknown.com:241/hello-world");

		assert.equal(URLToReadable("http://unknown.com/hello=world"), "//unknown.com/hello=world");
		assert.equal(URLToReadable("/hello=world"), "/hello=world");
		assert.equal(URLToReadable("hello=world"), "/hello=world");
	});

	exports.URLToLocation = URLToLocation;

	function URLToLocation(url) {
		if (typeof url === "undefined") {
			return {
				scheme: "",
				host: "",
				path: "",
				query: "",
				fragment: ""
			};
		}
		var rx = new RegExp("^((http|https):)?(//([^/?#]*))?([^?#]*)(\\?([^#]*))?(#(.*))?");
		var matches = url.match(rx);

		var path = matches[5] || "";
		if (path.charAt(0) !== "/") {
			path = "/" + path;
		}

		var fragment = matches[8] || "";
		if (fragment === "#") {
			fragment = "";
		}

		return {
			scheme: matches[2] || "",
			host: matches[4] || "",
			path: path,

			query: matches[6] || "",
			fragment: fragment || ""
		};
	}

	TestCase("URLToLocation", function(assert) {
		function verify(url) {
			var loc = URLToLocation(url);
			assert.equal(loc.scheme, "");
			assert.equal(loc.host, "");
			assert.equal(loc.query, "?q=csdfa asdf&filter=10.2.600");
			assert.equal(loc.fragment, "#alpha");
		}
		verify("/search=search?q=csdfa asdf&filter=10.2.600#alpha");
		verify("search=search?q=csdfa asdf&filter=10.2.600#alpha");
	});
});
