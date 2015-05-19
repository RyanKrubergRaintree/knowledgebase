'use strict';

import React from 'react';
import {convert} from "wiki/convert";

var Search = React.createClass({
	displayName: "Search",

	lastPageKey: undefined,
	search: function(ev){
		var Lineup = this.props.Lineup;
		var query = this.refs.query.getDOMNode().value.trim();
		this.lastPageKey = Lineup.open({
			url: '/index/search?q='+query,
			title: 'Search "' + query + '"',
			insteadOf: this.lastPageKey
		});
		ev.preventDefault();
	},
	keyDown: function(ev){
		var Lineup = this.props.Lineup;

		if(ev.keyCode == 27){// esc
			Lineup.closeLast();
			return;
		}
		if((ev.keyCode == 13) && (ev.ctrlKey || ev.shiftKey)){
			if(ev.ctrlKey){
				Lineup.clear();
			}
			var query = this.refs.query.getDOMNode().value.trim();
			Lineup.open({url: convert.LinkToURL(query), link:query});
			ev.preventDefault();
		}

		var pages = document.getElementsByClassName('page');
		if(pages.length == 0){
			return;
		};

		var page = pages[pages.length-1];
		var middle = page.getElementsByClassName('page-middle')[0];

		switch(ev.keyCode){
		case 33: // pageup
			middle.scrollTop -= middle.clientHeight
			break;
		case 34: // pagedown
			middle.scrollTop += middle.clientHeight
			break;
		}
	},
	render: function(){
		return React.DOM.form({
			className: 'page-search',
			onSubmit: this.search
		},
			React.DOM.input({
				ref:'query',
				placeholder: 'Search',
				onKeyDown: this.keyDown
			}, ''),
			React.DOM.input({type:'submit', value:'>'})
		);
	}
});

export var Header = React.createClass({
	displayName: 'Header',
	render: function(){
		var Lineup = this.props.Lineup;

		return React.DOM.header({id:'header'},
			React.DOM.span({},
				this.props.Title,
				' '
			),
			React.DOM.button({
				onClick: function(ev){
					if(ev.button != 1){
						Lineup.clear();
					}
					Lineup.open({url:'/home', link:'Home'});
				}
			}, 'Home'),
			' ',
			React.DOM.button({
				onClick: function(ev){
					Lineup.open({url:'/', link:''});
				}
			}, 'New Page'),
			React.DOM.button({
				onClick: function(ev){
					if(ev.button != 1){
						Lineup.clear();
					}
					Lineup.open({url:'/', link:''});
				}
			}, 'Clear'),
			' ',
			React.createElement(Search, {
				Lineup: this.props.Lineup
			})
		);
	}
});
