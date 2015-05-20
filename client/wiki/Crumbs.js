// import "/wiki/Lineup.js"
// import "/wiki/Wiki.js"
// import "/wiki/Convert.js"

(function(Wiki){
	"use strict";

	var separator = '| ';

	// converts to
	// #┃ref1┃ref2┃ref3
	function toHash(proxies){
		return separator + proxies.map(function(ref){
			var loc = Convert.URLToLocation(ref.url);
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

		var proxies = [];
		hash.split(separator).map(function(token){
			token = token.trim();
			if(token.trim() === ''){
				return;
			}

			var url = token;
			proxies.push(new Wiki.PageProxy({
				url: url,
				link: Convert.URLToLink(url),
				title: Convert.URLToTitle(url),
				key: -1
			}));
		});
		return proxies;
	}

	Wiki.Crumbs = Crumbs;
	function Crumbs(lineup){
		this.lineup_ = lineup;
		this.navigatingTo_ = "";

		this.lineup_.on('changed', this.lineupChanged, this);
		var self = this;
		window.addEventListener('hashchange', function(ev){
			if(window.location.hash !== self.navigatingTo_){
				self.lineup_.updateRefs(fromHash(window.location.hash));
			}
		})
	}

	Crumbs.prototype = {
		lineupChanged: function(event){
			this.navigatingTo_ = toHash(this.lineup_.proxies);
			window.location.hash = this.navigatingTo_;
		},
		initLineup: function(){
			var hash = window.location.hash;
			if((hash === '') || (hash === '#')){
				this.lineup_.open(Global.HomePage);
			} else {
				this.lineup_.updateRefs(fromHash(window.location.hash));
			}
		}
	};
})(Wiki);
