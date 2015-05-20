import React from "react";

import {Pages} from 'view/Pages';
import {Header} from 'view/Header';

import {global} from 'global';

export function init() {
	global.History.updateFromURL();
}

var App = React.createClass({
	displayName: "App",
	render: function () {
		return (
			React.DOM.div({ id: "root" },
				React.createElement(Pages,  global),
				React.createElement(Header, global)
			)
		)
	}
});

React.initializeTouchEvents(true);
React.render(
	React.createElement(App),
	document.getElementById("app")
);
