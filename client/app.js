package("kb.app", function(exports) {
	"use strict";

	window.KBUser = "Egon Elbre";
	window.KBHomePage = "Community=Welcome";

	depends("app.css");

	depends("Crumbs.js");
	depends("Lineup.js");
	depends("Site.js");
	depends("Selection.js");

	var app = exports;

	app.Lineup = new kb.Lineup();
	app.Crumbs = new kb.Crumbs(app.Lineup);
	app.CurrentSelection = new kb.Selection();

	function initialize(mountNode) {
		var site = React.createElement(kb.Site, {
			Lineup: app.Lineup,
			CurrentSelection: app.CurrentSelection
		});
		ReactDOM.render(site, mountNode);

		app.Crumbs.initLineup(KBHomePage);
		document.onclick = function(ev) {
			ev = ev || window.event; // IE8 passes the event in global scope
			app.Lineup.handleClickLink(ev);
		};
	}

	initialize(document.getElementById("site"));

	// closing of the last page
	document.onkeydown = function(ev) {
		ev = ev || window.event; // IE8 passes the event in global scope

		function elementIsEditable(elem) {
			return elem && (
				((elem.nodeName === "INPUT") && (elem.type === "text")) ||
				(elem.nodeName === "TEXTAREA") ||
				(elem.contentEditable === "true")
			);
		}

		if (ev.defaultPrevented || elementIsEditable(ev.target)) {
			return;
		}
		if (ev.keyCode === 27) {
			app.Lineup.closeLast();
			ev.preventDefault();
			ev.stopPropagation();
		}
	};
});
