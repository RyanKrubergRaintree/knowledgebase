'use strict';

import {PageRef} from './Lineup';
import {convert} from './convert';

var separator = '┃';

// converts to
// #┃ref1┃ref2┃ref3
function toHash(pagerefs){
	return separator + pagerefs.map(function(ref){
		var loc = convert.URLToLocation(ref.url);
		if((loc.origin == "") || (loc.origin == window.location.origin)) {
			return loc.pathname + loc.search;
		}
		return "//" + loc.host + loc.pathname + loc.search;
	}).join(separator);
}

// converts from
// #┃ref1┃ref2┃ref3
function fromHash(hash){
	// IE doesn't have the first #
	if(hash[0] == '#'){
		hash = hash.substr(2);
	} else {
		hash = hash.substr(1);
	}

	var pagerefs = [];
	hash.split(separator).map(function(token){
		token = token.trim();
		if(token.trim() === ''){
			return;
		}

		var url = token;
		pagerefs.push(new PageRef({
			url: url,
			link: convert.URLToLink(url),
			title: convert.URLToTitle(url),
			key: -1
		}));
	});
	return pagerefs;
}

export class History {
	constructor(lineup, defaultPage){
		this.lineup = lineup;
		this.defaultPage = defaultPage;

		var self = this;
		window.addEventListener('hashchange', function(ev){
			//TODO: review this on IE/Chrome
			// var location = convert.URLToLocation(ev.newURL);
			if(window.location.hash !== self.navigatingTo) {
				self.lineup.updateRefs(fromHash(window.location.hash));
			}
		});

		this.lineup.listen(function(pagerefs){
			self.navigatingTo = toHash(pagerefs);
			window.location.hash = self.navigatingTo;
		});
	}
	updateFromURL(){
		var hash = window.location.hash;
		if(hash === ''){
			this.lineup.open(this.defaultPage);
		} else {
			this.lineup.updateRefs(fromHash(window.location.hash));
		}
	}
}
