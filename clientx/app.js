//import "/global/polyfills.js"
//import "/view/App.js"

Global = {
	User: ""
};

function initialize(mountNode){
	React.initializeTouchEvents(true);
	React.render(React.createElement(View.App, null), mountNode);
}
