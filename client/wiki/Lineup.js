// import "/util/Notifier.js"
// import "/wiki/Wiki.js"
// import "/wiki/Convert.js"

(function(Wiki){
	"use strict";

	Wiki.Lineup = Lineup;
	function Lineup(){
		this.pagerefs = [];
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
			for(var i = 0; i < this.pagerefs.length; i += 1){
				if(this.pagerefs[i].key == key){
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
				this.pagerefs = this.pagerefs.slice(0, i + 1);
			}
		},

		clear: function(){
			this.pagerefs = [];
			this.changed();
		},

		close: function(key){
			this.pagerefs = this.pagerefs.filter(function(pageref){
				return pageref.key !== key;
			});
			this.changed();
		},

		closeLast: function(){
			var pagerefs = this.pagerefs;
			pagerefs = pagerefs.slice(0, Math.max(pagerefs.length-1, 1));
			this.changed();
		},

		changeRef: function(key, page){
			var i = this.indexOf_(key);
			if(i >= 0){
				var ref = this.pagerefs[i];
				ref.url = Convert.URLToReadable(page.url);
				ref.link = Convert.URLToLink(page.link);
				ref.title = page.title;
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

			var pageref = new PageRef({
				url: url,
				title: props.title || Convert.URLToTitle(url),
				link: props.link || Convert.URLToLink(url),
				key: this.lastKey++
			});

			if(props.link === ""){
				pageref.link = "";
			}

			var i = this.indexOf_(props.insteadOf);
			if(i >= 0){
				this.pagerefs[i] = pageref;
			} else {
				this.pagerefs.push(pageref);
			}

			this.changed();
			return pageref.key;
		},


		updateRefs: function(nextrefs){
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
	};

	Wiki.PageRef = PageRef;
	function PageRef(props) {
		this.url = props.url;
		this.owner = props.owner;
		this.link = props.link;
		this.title = props.title;
		this.key = props.key;
	}
})(Wiki);
