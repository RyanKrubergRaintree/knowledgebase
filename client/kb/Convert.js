package('kb.convert', function(exports){
	'use strict';

	depends('unicode/identifier.js');
	depends('unicode/runename.js');

	function trimProtocol(link){
		if(link.indexOf('https://', link) === 0){
			return link.substr(6).trim();
		} else if(link.indexOf('http://', link) === 0){
			return link.substr(5).trim();
		}
		return link.trim();
	}

	function trimLeadingSlashes(link){
		// remove prefix '/'
		while(link[0] === '/') {
			link = link.substr(1);
		}
		return link;
	}


	// TextToSlug converts text to a slug
	//
	// * numbers and '/' are left intact
	// * letters will be lowercased (if possible)
	// * '-', ',', '.', ' ', '_' will be converted to '-'
	// * other symbols or punctuations will be converted to html entity reference name
	//   (if there exists such reference name)
	// * everything else will be converted to '-'
	//
	// Example:
	//   '&Hello_世界/+!' ==> 'amp-hello-世界/plus-excl'
	//   'Hello  World  /  Test' ==> 'hello-world/test'
	exports.TextToSlug = TextToSlug;
	function TextToSlug(title){
		var cutdash = true,
			emitdash = false;

		var slug = '';

		for(var i = 0; i < title.length; i += 1){
			var r = title.charAt(i);
			if(kb.unicode.IsIdent(r)) {
				if(emitdash && !cutdash){
					slug += '-';
				}
				slug += r.toLowerCase();

				emitdash = false;
				cutdash = false;
				continue;
			}
			if((r === '/') || (r === ':')){
				slug += r;
				emitdash = false;
				cutdash = true;
			} else if ((r === '-') || (r === ',') || (r === '.') || (r === ' ') || (r === '_')) {
				emitdash = true;
			} else {
				var name = kb.unicode.RuneName[r];
				if(name){
					if(!cutdash){
						slug += '-';
					}
					slug += name;
					cutdash = false;
				}
				emitdash = true;
			}
		}

		if(slug.length === 0){
			return '-';
		}

		return slug;
	}

	// There are several possible links
	// 'http://kb.example.com/example'
	// 'https://kb.example.com/example'
	// '//kb.example.com/example'
	// '/kb:example' - rooted local URL
	// 'kb:Example' - local URL
	exports.LinkToReference = LinkToReference;
	function LinkToReference(link){
		link = trimProtocol(link);
		// External site:
		// '//kb.example.com/example'
		if((link[0] === '/') && (link[1] === '/') ) {
			return {
				link: URLToReadable(link),
				url:  link,
				title: LinkToTitle(link)
			};
		}

		link = trimLeadingSlashes(link);
		var i = link.indexOf(':');
		var owner = i >= 0 ? link.substr(0,i): '';

		return {
			link: URLToReadable(link),
			owner: owner,
			url: '/' + TextToSlug(link),
			title: LinkToTitle(link)
		};
	}

	exports.ReferenceToLink = ReferenceToLink;
	function ReferenceToLink(ref){
		return URLToReadable(ref.url);
	}

	exports.LinkToTitle = LinkToTitle;
	function LinkToTitle(link){
		link = trimProtocol(link);
		link = trimLeadingSlashes(link);
		var i = Math.max(link.lastIndexOf('/'), link.indexOf(':'));
		link = link.substr(i + 1);
		return link;
	}

	exports.LinkToOwner = LinkToOwner;
	function LinkToOwner(link){
		link = URLToReadable(link);
		link = trimProtocol(link);
		link = trimLeadingSlashes(link);

		var i = link.indexOf(':');
		link = link.substr(0, i);
		return link.trim().toLowerCase();
	}

	exports.URLToReadable = URLToReadable;
	function URLToReadable(url){
		var loc = URLToLocation(url);
		if((typeof loc.origin === 'undefined') || (loc.origin === window.location.origin)){
			if(loc.pathname[0] === '/'){
				return loc.pathname + loc.search + loc.hash;
			}
			return '/' + loc.pathname + loc.search + loc.hash;
		}
		return url;
	}

	exports.URLToLocation = URLToLocation;
	function URLToLocation(url){
		var a = document.createElement('a');
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
			href: a.href
		};
	}
});