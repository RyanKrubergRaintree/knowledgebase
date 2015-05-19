'use strict';

import {Page} from './Page';
import {Notifier} from 'util/Notifier';
import {convert} from './convert';


function jsonurl(url){
	var loc = convert.URLToLocation(url)
	//TODO: check whether pathname already contains .json
	return "//" + loc.host + loc.pathname + ".json" + loc.hash + loc.search;
}

function stripjson(url){
	//TODO: strip properly
	return url.replace(".json", "");
}

export class PageProxy {
	constructor(pageref, store, lineup){
		this.store = store;
		this.lineup = lineup;

		this.url = pageref.url;
		this.pageref = pageref;
		this.patching = false;
		this.pending = [];
		this.callbacks = [];
		this.page = null;

		this.state = "loading";
		this.status = 200;
		this.statusText = "";
		this.responseText = "";
	}

	initIfNeeded(){
		if(this.page === null){
			this.page = new Page({
				url: this.pageref.url,
				title: this.pageref.title,
				flags: ['read-only']
			});
			this.changed();
		}
	}

	listen(onchange){
		this.callbacks.push(onchange);

		var self = this;
		return function(){
			self.unlisten(onchange);
		};
	}

	unlisten(onchange){
		this.callbacks = this.callbacks.filter(function(cb){
			return cb != onchange;
		});

		if(this.callbacks.length == 0){
			this.store.removeProxy(this);
		}
	}

	changed(){
		var self = this;
		this.callbacks.map(function(cb){
			cb(self);
		});
	}

	// returns true if it's a success status
	updateStatus(xhr){
		this.state = "loaded";

		if((xhr.status >= 300) || (xhr.status < 200)){
			this.state = "errored";

			if(xhr.status === 404){
				this.state = "not-found";
			}

			this.status = xhr.status;
			this.statusText = xhr.statusText;
			this.responseText = xhr.responseText;

			return false;
		}
		return true;
	}

	// CALLBACKS FOR XHR
	loaded(ev){
		var xhr = ev.target;
		if(!this.updateStatus(xhr)){
			this.changed();
			return;
		}

		var data = JSON.parse(xhr.response),
			page = new Page(data);
		if(xhr.responseURL){
			page.updateURL(stripjson(xhr.responseURL));
		}
		this.page = page;
		this.state = "loaded";
		this.changed();
	}

	errored(ev){
		var xhr = ev.target;

		this.status = "";
		this.statusText = "";
		this.responseText = "";

		if(typeof xhr !== 'undefined'){
			this.updateStatus(xhr);
		} else {
			this.status = "failed";
			this.message = ev;
		}
		this.state = "failed";
	}

	deleted(ev){
		this.reload();
	}

	created(ev){
		this.loaded(ev);
		this.lineup.changeRef(this.pageref.key, this.page);
	}

	patched(ev){
		this.patching = false;
		var xhr = ev.target;
		// failed to patch
		if(!this.updateStatus(xhr)){
			this.pending = [];
			this.reload();
			this.changed();
			return;
		}
		this._sendpatch();
	}
	patcherrored(ev){
		var xhr = ev.target;
		this.pending = [];
		this.patching = false;
		this.reload();
	}

	// API
	create(props){
		var url = convert.LinkToURL(props.link);
		var title = convert.LinkToTitle(props.link);

		var xhr = new XMLHttpRequest();
		xhr.withCredentials = true;
		xhr.addEventListener('load', this.created.bind(this), false);
		xhr.addEventListener('error', this.errored.bind(this), false);

		xhr.open('PUT', jsonurl(url), true);

		xhr.setRequestHeader('Accept', 'application/json');
		xhr.setRequestHeader('Content-Type', 'application/json');

		xhr.send(JSON.stringify({
			title: title,
			tags: props.tags,
			synopsis: props.synopsis
		}));
	}

	reload(){
		var xhr = new XMLHttpRequest();
		xhr.withCredentials = true;
		xhr.addEventListener('load', this.loaded.bind(this), false);
		xhr.addEventListener('error', this.errored.bind(this), false);

		xhr.open('GET', jsonurl(this.url), true);
		xhr.setRequestHeader('Accept', 'application/json');
		xhr.send();
	}

	_sendpatch(){
		var patch = this.pending.shift();
		if(patch){
			this.patching = true;

			var xhr = new XMLHttpRequest();
			xhr.withCredentials = true;
			xhr.addEventListener('load', this.patched.bind(this), false);
			xhr.addEventListener('error', this.patcherrored.bind(this), false);

			xhr.open('PATCH', jsonurl(this.url), true);

			xhr.setRequestHeader('Accept', 'application/json');
			xhr.setRequestHeader('Content-Type', 'application/json');

			xhr.send(JSON.stringify(patch));
		}
	}

	modify(op){
		var page = this.page;
		var version = page.version;
		this.page.apply(op);
		this.pending.push(op);

		this._sendpatch();
		this.changed()
	}

	delete(){
		var xhr = new XMLHttpRequest();
		xhr.withCredentials = true;
		xhr.open('DELETE', jsonurl(this.url), true);

		xhr.addEventListener('load', this.deleted.bind(this), false);
		xhr.addEventListener('error', this.errored.bind(this), false);

		xhr.setRequestHeader('Accept', 'application/json');
		xhr.setRequestHeader('Content-Type', 'application/json');

		xhr.send();
		return xhr;
	}
}
