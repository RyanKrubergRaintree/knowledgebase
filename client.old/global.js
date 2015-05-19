'use strict';

import {Lineup} from 'wiki/Lineup';
import {History} from 'wiki/History';
import {PageStore} from 'wiki/PageStore';
import {convert} from "wiki/convert";

var StartingPage = {
	url: "/home",
	title: "Home",
	link: "Home"
};

export var global = {
	Title: 'KnowledgeBase',

	Lineup: null,
	History: null,
	Store: null
};

global.Lineup = new Lineup();
global.History = new History(global.Lineup, StartingPage);
global.Store = new PageStore(global.Lineup);

function localURL(url) {
	return !((url[0] == "/") && (url[1] == "/"));
}

export function OpenNewPage(ev){
	function findPageNode(el){
		for(var i = 0; i < 32; i += 1){
			if(el == null){ return null; }
			if(el.classList.contains("page")){
				return el;
			}
			el = el.parentElement;
		}
		return undefined;
	}


	var target = ev.target;
	var pagenode = findPageNode(target);

	var url = target.href;
	if(pagenode) {
		var locFrom = convert.URLToLocation(pagenode.dataset.url);
		var locTo = convert.URLToLocation(url);
		url = "//" + locFrom.host + locTo.pathname;
	}

	var link = target.dataset.link;
	var title = target.innerText;

	if(ev.button == 1){
		global.Lineup.open({
			url: url,
			link: link,
			title: title
		});
	} else {
		var key = undefined;
		if(pagenode){
			key = parseInt(pagenode.dataset.key);
		}
		global.Lineup.open({
			url: url,
			link: link,
			title: title,
			after: key
		});
	}
	ev.preventDefault();
}

window.addEventListener("click", function(ev){
	if(ev.target.localName != "a") return;
	if(ev.target.classList.contains("external-link")) return;
	if(ev.target.onclick != null) return;
	if(ev.target.onmousedown != null) return;
	if(ev.target.onmouseup != null) return;
	if(ev.target.href == "") return;
	OpenNewPage(ev);
});


function elementIsEditable(elem){
	return elem && (
		((elem.nodeName === 'INPUT') && (elem.type === 'text')) ||
		(elem.nodeName === 'TEXTAREA') ||
		(elem.contentEditable === 'true')
 	);
}

// closing of the last page
window.addEventListener("keydown", function(ev){
	if(ev.defaultPrevented || elementIsEditable(ev.target)){
		return;
	}
	if(ev.keyCode == 27){
		global.Lineup.closeLast();
	}
});
