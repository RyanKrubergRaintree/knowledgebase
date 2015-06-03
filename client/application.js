import "global.js";
import "crumbs.js";
import "lineup.js";
import "site.js";

Global.Lineup = new KB.Lineup();
Global.Crumbs = new KB.Crumbs(Global.Lineup);

function initialize(mountNode){
	React.initializeTouchEvents(true);
	var site = React.createElement(KB.Site, {Lineup: Global.Lineup});
	React.render(site, mountNode);

	Global.Crumbs.initLineup(Global.HomePage);
	window.addEventListener("click",
		Global.Lineup.handleClickLink.bind(Global.Lineup));
}

initialize(document.getElementById("site"));

// closing of the last page
window.addEventListener("keydown", function(ev){
	function elementIsEditable(elem){
		return elem && (
			((elem.nodeName === 'INPUT') && (elem.type === 'text')) ||
			(elem.nodeName === 'TEXTAREA') ||
			(elem.contentEditable === 'true')
	 	);
	}

	if(ev.defaultPrevented || elementIsEditable(ev.target)){
		return;
	}
	if(ev.keyCode == 27){
		Global.Lineup.closeLast();
		ev.preventDefault();
		ev.stopPropagation();
	}
});
