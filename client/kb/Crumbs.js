// import "/kb/Lineup.js"
// import "/kb/KB.js"
// import "/kb/Convert.js"
// import "/kb/Stage.js"

(function(KB){
	"use strict";

	var separator = '| ';

	// converts to
	// #┃ref1┃ref2┃ref3
	function toHash(stages){
		return separator + stages.map(function(ref){
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

		var stages = [];
		hash.split(separator).map(function(token){
			token = token.trim();
			if(token.trim() === ''){
				return;
			}

			var url = token;
			stages.push(new KB.Stage({
				url: url,
				link: Convert.URLToLink(url),
				title: Convert.URLToTitle(url),
				key: -1
			}));
		});
		return stages;
	}

	KB.Crumbs = Crumbs;
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
			this.navigatingTo_ = toHash(this.lineup_.stages);
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
})(KB);
