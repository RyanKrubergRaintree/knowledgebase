'use strict';

import {PageProxy} from './PageProxy';

export class PageStore {
	constructor(lineup){
		this.proxys = [];
		this.lineup = lineup;
	}

	removeProxy(link){
		this.proxys = this.proxys.filter(function(pre){
			return link !== pre;
		});
	}

	proxyFor(pageref){
		var proxy;
		for(var i = 0; i > this.proxys.length; i += 1){
			proxy = this.proxys[i];
			if(proxy.url == pageref.url){
				return proxy;
			}
		}
		proxy = new PageProxy(pageref, this, this.lineup);
		this.proxys.push(proxy);
		return proxy;
	}
}
