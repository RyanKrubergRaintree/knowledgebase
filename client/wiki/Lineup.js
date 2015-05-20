// import "/util/Notifier.js"
// import "/wiki/Wiki.js"
// import "/wiki/Convert.js"
// import "/wiki/PageProxy.js"

(function(Wiki){
	"use strict";

	Wiki.Lineup = Lineup;
	function Lineup(){
		this.proxies = [];
		this.lastKey = 0;
		this.notifier = new Notifier();
	}

	Lineup.prototype = {
		on: function(event, handler, recv){ this.notifier.on(event, handler, recv); },
		off: function(event, handler, recv){ this.notifier.off(event, handler, recv); },
		remove: function(recv){ this.notifier.remove(recv); },

		changed: function(){
			this.notifier.emit({
				type:"changed",
				lineup: this
			});
		},

		indexOf_: function(key){
			if(typeof key === 'undefined'){
				return -1;
			}
			for(var i = 0; i < this.proxies.length; i += 1){
				if(this.proxies[i].key == key){
					return i;
				}
			}
			return -1;
		},

		trim_: function(key){
			if(typeof key === 'undefined'){
				return;
			}
			var i = this.indexOf_(key);
			if(i >= 0){
				this.proxies = this.proxies.slice(0, i + 1);
			}
		},

		clear: function(){
			this.proxies = [];
			this.changed();
		},

		close: function(key){
			this.proxies = this.proxies.filter(function(proxy){
				return proxy.key !== key;
			});
			this.changed();
		},

		closeLast: function(){
			var proxies = this.proxies;
			proxies = proxies.slice(0, Math.max(proxies.length-1, 1));
			this.changed();
		},

		changeRef: function(key, proxy){
			var i = this.indexOf_(key);
			if(i >= 0){
				var ref = this.proxies[i];
				ref.url = Convert.URLToReadable(proxy.url);
				ref.link = Convert.URLToLink(proxy.link);
				ref.title = proxy.title;
				this.changed();
			}
		},


		// url
		// title, optional
		// link, optional
		// after, optional
		// insteadOf, optional
		open: function(props){
			this.trim_(props.after);
			var url = Convert.URLToReadable(props.url);

			var proxy = new Wiki.PageProxy({
				url: url,
				title: props.title || Convert.URLToTitle(url),
				link: props.link || Convert.URLToLink(url),
				key: this.lastKey++
			});

			if(props.link === ""){
				proxy.link = "";
			}

			var i = this.indexOf_(props.insteadOf);
			if(i >= 0){
				this.proxies[i] = proxy;
			} else {
				this.proxies.push(proxy);
			}

			this.changed();
			return proxy.key;
		},


		updateRefs: function(nextproxies){
			var preproxies = this.proxies.slice();
			var changed = false;

			var self = this;
			var newproxies = nextproxies.map(function(proxy){
				var prev = preproxies.shift();
				if(prev && (prev.url == proxy.url)){
					return prev;
				}
				changed = true;
				proxy.key = self.lastKey++;
				return proxy;
			});

			if(preproxies.length > 0){
				changed = true;
			}
			if(changed){
				this.proxies = newproxies;
				this.changed();
			}
		}
	};
})(Wiki);
