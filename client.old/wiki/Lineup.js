'use strict';

import {convert} from "./convert";
import {Notifier} from "util/Notifier";

export class PageRef {
	constructor(props){
		this.url = props.url;
		this.link = props.link;
		this.title = props.title;
		this.key = props.key;
	}
}

export class Lineup extends Notifier {
	constructor(){
		super();

		this.pagerefs = [];
		this.lastKey = 0;

		var self = this;
		this.boundHashChanged = function(ev){self.hashchanged(ev)};
	}

	changed(){
		this.notify(this.pagerefs);
	}

	_indexOf(key){
		if(typeof key === 'undefined'){
			return -1;
		}
		for(var i = 0; i < this.pagerefs.length; i += 1){
			if(this.pagerefs[i].key == key){
				return i;
			}
		}
		return -1;
	}

	_trim(key){
		if(typeof key === 'undefined'){
			return;
		}
		var i = this._indexOf(key);
		if(i >= 0){
			this.pagerefs = this.pagerefs.slice(0, i + 1);
		}
	}

	clear(){
		this.pagerefs = [];
		this.changed();
	}

	close(key){
		this.pagerefs = this.pagerefs.filter(function(pageref){
			return pageref.key !== key;
		});
		this.changed();
	}

	closeLast(){
		var pagerefs = this.pagerefs;
		pagerefs = pagerefs.slice(0, Math.max(pagerefs.length-1, 1));
		this.changed();
	}

	changeRef(key, page){
		var i = this._indexOf(key);
		if(i >= 0){
			var ref = this.pagerefs[i];
			ref.url = convert.URLToReadable(page.url);
			ref.link = convert.URLToLink(page.link);
			ref.title = page.title;
			this.changed();
		}
	}

	// url
	// title, optional
	// link, optional
	// after, optional
	// insteadOf, optional
	open(props){
		this._trim(props.after);
		var url = convert.URLToReadable(props.url);
		var pageref = new PageRef({
			url: url,
			title: props.title || convert.URLToTitle(url),
			link: props.link || convert.URLToLink(url),
			key: this.lastKey++
		});

		if(props.link === ""){
			pageref.link = "";
		}

		var i = this._indexOf(props.insteadOf);
		if(i >= 0){
			this.pagerefs[i] = pageref;
		} else {
			this.pagerefs.push(pageref);
		}

		this.changed();
		return pageref.key;
	}

	updateRefs(nextrefs){
		var prerefs = this.pagerefs.slice();
		var changed = false;

		var self = this;
		var newrefs = nextrefs.map(function(pageref){
			var prev = prerefs.shift();
			if(prev && (prev.url == pageref.url)){
				return prev;
			}
			changed = true;
			pageref.key = self.lastKey++;
			return pageref;
		});

		if(prerefs.length > 0){
			changed = true;
		}
		if(changed){
			this.pagerefs = newrefs;
			this.changed();
		}
	}
}
