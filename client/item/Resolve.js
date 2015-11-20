package("kb.item", function(exports) {
	"use strict";

	depends("../Convert.js");

	var rxExternalLink = /\[\[\s*(https?\:[^ \]]+)\s+([^\]]+)\]\]/g;
	var rxInternalLink = /\[\[\s*([^\]]+)\s*\]\]/g;

	exports.Resolve = Resolve;

	function Resolve(stage, text) {
		if ((typeof text === "undefined") || (text === null)) {
			return "";
		}

		text = text.replace(rxExternalLink,
			"<a href=\"$1\" class=\"external-link\" target=\"_blank<div>\"</div> rel=\"nofollow\">$2</a>");

		text = text.replace(rxInternalLink, function(match, link) {
			var ref = kb.convert.LinkToReference(link, stage);
			return "<a href=\"" + ref.url + "\" data-link=\"" + ref.link + "\" >" + ref.title + "</a>";
		});
		return text;
	}

	//TODO: add HTML sanitation
	exports.ResolveHTML = ResolveHTML;

	function ResolveHTML(stage, text) {
		return Resolve(stage, text);
	}
});
