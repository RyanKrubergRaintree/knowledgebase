package("kb", function(exports) {
	"use strict";

	depends("Convert.js");

	// converts to
	// #┃ref1┃ref2┃ref3
	function stagesToHash(stages) {
		return "#" + stages.map(kb.convert.ReferenceToLink).join("|");
	}

	// converts from
	// #┃ref1┃ref2┃ref3
	function hashToRefs(hash) {
		// IE doesn"t have the first #
		if (hash[0] === "#") {
			hash = hash.substr(1);
		} else {
			hash = hash.substr(0);
		}

		var refs = [];
		hash.split("|").map(function(link) {
			if (link.trim() === "") {
				return;
			}
			refs.push(kb.convert.LinkToReference(link));
		});
		return refs;
	}

	exports.Crumbs = Crumbs;

	function Crumbs(lineup) {
		this.lineup_ = lineup;
		this.navigatingTo_ = "";
	}

	Crumbs.prototype = {
		attach: function(defaultPages, home) {
			var hash = window.location.hash;
			if ((hash === "") || (hash === "#")) {
				if (defaultPages.length > 0) {
					this.lineup_.openPages(defaultPages);
				} else {
					this.lineup_.openLink(home);
				}
			} else {
				this.lineup_.updateRefs(hashToRefs(window.location.hash));
			}
			this.lineup_.on("changed", this.lineupChanged, this);

			var self = this;
			window.onhashchange = function( /*ev*/ ) {
				if (window.location.hash !== self.navigatingTo_) {
					self.lineup_.updateRefs(hashToRefs(window.location.hash));
				}
			};
		},
		detach: function() {
			window.onhashchange = null;
			this.lineup_.remove(this);
		},

		lineupChanged: function( /*ev*/ ) {
			this.navigatingTo_ = stagesToHash(this.lineup_.stages);
			window.location.hash = this.navigatingTo_;
		}
	};
});
