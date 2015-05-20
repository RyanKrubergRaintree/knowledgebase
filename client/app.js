// import "/global/polyfills.js"

Global = {
	User: "",
	HomePage: {},

	Lineup: null,
	Crumbs: null
};

// import "/wiki/Crumbs.js"
// import "/wiki/Lineup.js"

Global.Lineup = new Wiki.Lineup();
Global.Crumbs = new Wiki.Crumbs(Global.Lineup);

// import "/view/App.js"
function initialize(mountNode){
	React.initializeTouchEvents(true);
	React.render(React.createElement(View.App, null), mountNode);

	Global.Crumbs.initLineup();
}
