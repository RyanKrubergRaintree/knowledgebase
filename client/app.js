// import "/global/polyfills.js"

Global = {
	User: "",
	HomePage: "",

	Lineup: null,
	Crumbs: null,
};

// import "/kb/Crumbs.js"
// import "/kb/Lineup.js"

Global.Lineup = new KB.Lineup();
Global.Crumbs = new KB.Crumbs(Global.Lineup);

// import "/kb/Site.js"
function initialize(mountNode){
	React.initializeTouchEvents(true);
	React.render(React.createElement(KB.Site, {
		Lineup: Global.Lineup
	}), mountNode);

	Global.Crumbs.initLineup(Global.HomePage);
	window.addEventListener("click",
		Global.Lineup.handleClickLink.bind(Global.Lineup));
}
