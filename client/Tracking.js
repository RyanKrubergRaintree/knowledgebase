package("kb", function (exports) {
	"use strict";

	exports.TrackPageView = function (url, title) {
		if (typeof ga !== "undefined") {
			ga("set", {
				page: url,
				title: title
			});
			ga("send", "pageview");
		}
	};
});
