import "lineup.js";
import "kb.js";
import "convert.js";
import "stage.js";

KB.Crumbs = (function(){
	"use strict";

	// converts to
	// #┃ref1┃ref2┃ref3
	function toHash(stages){
		return "# " + stages.map(Convert.ReferenceToLink).join('| ');
	}

	// converts from
	// #┃ref1┃ref2┃ref3
	function fromHash(hash){
		// IE doesn't have the first #
		if(hash[0] == '#'){
			hash = hash.substr(1);
		} else {
			hash = hash.substr(0);
		}

		var stages = [];
		hash.split('|').map(function(link){
			if(link.trim() === ''){ return; }
			stages.push(new KB.Stage(Convert.LinkToReference(link)));
		});
		return stages;
	}

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
		initLineup: function(defaultPage){
			var hash = window.location.hash;
			if((hash === '') || (hash === '#')){
				this.lineup_.openLink(defaultPage);
			} else {
				this.lineup_.updateRefs(fromHash(window.location.hash));
			}
		}
	};

	return Crumbs;
})(KB);
