package('kb.app', function(exports){
	'use strict';

	depends('app.css');

	depends('Crumbs.js');
	depends('Lineup.js');
	depends('Site.js');

	var app = exports;

	app.Lineup = new kb.Lineup();
	app.Crumbs = new kb.Crumbs(app.Lineup);

	function initialize(mountNode){
		React.initializeTouchEvents(true);
		var site = React.createElement(kb.Site, {Lineup: app.Lineup});
		React.render(site, mountNode);

		app.Crumbs.initLineup(KBHomePage);
		window.addEventListener('click',
			app.Lineup.handleClickLink.bind(app.Lineup));
	}

	initialize(document.getElementById('site'));

	// closing of the last page
	window.addEventListener('keydown', function(ev){
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
		if(ev.keyCode === 27){
			app.Lineup.closeLast();
			ev.preventDefault();
			ev.stopPropagation();
		}
	});
});