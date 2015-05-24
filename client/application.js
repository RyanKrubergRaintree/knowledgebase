import "global.js";
import "crumbs.js";
import "lineup.js";

Global.Lineup = new KB.Lineup();
Global.Crumbs = new KB.Crumbs(Global.Lineup);

import "site.js";
function initialize(mountNode){
	React.initializeTouchEvents(true);
	var site = React.createElement(KB.Site, {Lineup: Global.Lineup});
	React.render(site, mountNode);

	Global.Crumbs.initLineup(Global.HomePage);
	window.addEventListener("click",
		Global.Lineup.handleClickLink.bind(Global.Lineup));
}

initialize(document.getElementById("site"));
