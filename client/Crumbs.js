package("kb", function(exports) {
	"use strict";

	depends("Lineup.js");
	depends("Convert.js");
	depends("Stage.js");

	// converts to
	// #┃ref1┃ref2┃ref3
	function toHash(stages) {
		return "#" + stages.map(kb.convert.ReferenceToLink).join("|");
	}

	// converts from
	// #┃ref1┃ref2┃ref3
	function fromHash(hash) {
		// IE doesn"t have the first #
		if (hash[0] === "#") {
			hash = hash.substr(1);
		} else {
			hash = hash.substr(0);
		}

		var stages = [];
		hash.split("|").map(function(link) {
			if (link.trim() === "") {
				return;
			}
			stages.push(new kb.Stage(kb.convert.LinkToReference(link)));
		});
		return stages;
	}

	exports.Crumbs = Crumbs;

	function Crumbs(lineup) {
		this.lineup_ = lineup;
		this.navigatingTo_ = "";

		this.lineup_.on("changed", this.lineupChanged, this);
		var self = this;
		window.onhashchange = function( /*ev*/ ) {
			if (window.location.hash !== self.navigatingTo_) {
				self.lineup_.updateRefs(fromHash(window.location.hash));
			}
		};
	}

	Crumbs.prototype = {
		lineupChanged: function( /*ev*/ ) {
			this.navigatingTo_ = toHash(this.lineup_.stages);
			window.location.hash = this.navigatingTo_;
		},
		initLineup: function(defaultPage) {
			var hash = window.location.hash;
			if ((hash === "") || (hash === "#")) {
				this.lineup_.openLink(defaultPage);
			} else {
				this.lineup_.updateRefs(fromHash(window.location.hash));
			}
		}
	};
});
