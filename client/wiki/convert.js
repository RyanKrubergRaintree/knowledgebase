'use strict';

import {IsIdent} from "./internal/unicode.js";
import {RuneName} from "./internal/runename.js";

export var convert = {
	Slugify: Slugify,
	LinkToURL: LinkToURL,
	LinkToTitle: LinkToTitle,
	URLToReadable: URLToReadable,
	URLToLink: URLToLink,
	URLToTitle: URLToTitle,
	URLToLocation: URLToLocation,
	resolve: resolve
};

// Slugify converts text to a slug
//
// * numbers and '/' are left intact
// * letters will be lowercased (if possible)
// * '-', ',', '.', ' ', '_' will be converted to '-'
// * other symbols or punctuations will be converted to html entity reference name
//   (if there exists such reference name)
// * everything else will be converted to '-'
//
// Example:
//   "&Hello_世界/+!" ==> "amp-hello-世界/plus-excl"
//   "Hello  World  /  Test" ==> "hello-world/test"
export function Slugify(title){
	var cutdash = true,
		emitdash = false;

	var slug = '';

	for(var i = 0; i < title.length; i += 1){
		var r = title[i];
		if(IsIdent(r)) {
			if(emitdash && !cutdash){
				slug += '-';
			}
			slug += r.toLowerCase();

			emitdash = false;
			cutdash = false;
			continue;
		}
		if(r == '/'){
			slug += r;
			emitdash = false;
			cutdash = true;
		} else if ((r == '-') || (r == ',') || (r == '.') || (r == ' ') || (r == '_')) {
			emitdash = true;
		} else {
			var name = RuneName[r];
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

	if(slug.length == 0){
		return '-';
	}

	return slug;
}

// [[Welcome Visitors]]          -> /welcome-visitors
// [[question/Welcome Visitors]] -> /question/welcome-visitors
// [[/index/Welcome Visitors]]   -> /index/welcome-visitors
export function LinkToURL(link){
	if(link.length == 0){ return "/"; }
	link = Slugify(link);

	// not properly rooted, this means we want a page
	if(link[0] != "/"){
		return "/" + link;
	}
	return link;
}

export function LinkToTitle(link){
	var i = link.lastIndexOf("/");
	return link.substr(i + 1);
}

// returns shortened path, if it's from current domain
export function URLToReadable(url){
	var loc = URLToLocation(url);
	if((typeof loc.origin == "undefined") || (loc.origin == window.location.origin)){
		if(loc.pathname[0] == "/") {
			return loc.pathname + loc.search + loc.hash;
		}
		return "/" + loc.pathname + loc.search + loc.hash;
	}
	return url;
}

export function URLToLink(url){
	return URLToReadable(url);
}

//TODO: make a better function
export function URLToTitle(url){
	return LinkToTitle(url);
}

export function URLToLocation(url){
	var a = document.createElement("a");
	a.href = url;
	return {
		get hash(){ return a.hash; },
		set hash(v){ a.hash = v; },
		get search(){ return a.search; },
		set search(v){ a.search = v; },
		get pathname(){ return a.pathname; },
		set pathname(v){ a.pathname = v; },
		get port(){ return a.port; },
		set port(v){ a.port = v; },
		get hostname(){ return a.hostname; },
		set hostname(v){ a.hostname = v; },
		get host(){ return a.host; },
		set host(v){ a.host = v; },
		get password(){ return a.password; },
		set password(v){ a.password = v; },
		get username(){ return a.username; },
		set username(v){ a.username = v; },
		get protocol(){ return a.protocol; },
		set protocol(v){ a.protocol = v; },
		get origin(){ return a.origin; },
		set origin(v){ a.origin = v; },
		get url(){ return a.href; },
		set url(v){ a.href = v; }
	};
}

var rxExternalLink = /\[\[\s*(http[a-zA-Z0-9:\.\/\-]+)\s+([^\]]+)\]\]/g,
	rxInternalLink = /\[\[\s*([^\]]+)\s*\]\]/g,
	externalLinkGlyph = '<span class="oi" data-glyph="external-link" />';

export function resolve(text){
	if(text == null) {
		return '';
	}

	// text = text.replace("&", "&amp;");
	// text = text.replace("<", "&lt;");
	// text = text.replace(">", "&gt;");
	text = text.replace(rxExternalLink,
		'<a href="$1" class="external-link" target="_blank" rel="nofollow">$2' + externalLinkGlyph + '</a>');

	text = text.replace(rxInternalLink, function(match, link){
		var url = LinkToURL(link),
			title = LinkToTitle(link);
		return '<a href="' + url + '" data-link="' + link + '" >' + title + "</a>";
	});

	return text;
}
