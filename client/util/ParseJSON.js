package("kb.util", function (exports) {
	"use strict";

	exports.ParseJSON = ParseJSON;

	function ParseJSON(data) {
		try {
			var result = JSON.parse(data);
		} catch (err) {
			console.error("Parsing failed:", data);
			throw err;
		}
		return result;
	}
});
