// import "/global/polyfills.js"

Global = {
	User: "",
	HomePage: {},

	Lineup: null,
	Crumbs: null
};

// import "/kb/Crumbs.js"
// import "/kb/Lineup.js"

Global.Lineup = new KB.Lineup();
Global.Crumbs = new KB.Crumbs(Global.Lineup);

// import "/view/App.js"
function initialize(mountNode){
	React.initializeTouchEvents(true);
	React.render(React.createElement(View.App, null), mountNode);

	Global.Crumbs.initLineup();
}
