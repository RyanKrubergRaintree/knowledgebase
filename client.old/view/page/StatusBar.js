'use strict';

import React from 'react';
import {convert} from 'wiki/convert';
import {newIcon} from 'util/icon';

export var StatusBar = React.createClass({
	displayName: 'StatusBar',

	getInitialState: function(){
		return {
			noteHidden: false
		};
	},

	hideNote: function(){ this.setState({noteHidden: true}); },
	showNote: function(){ this.setState({noteHidden: false}); },

	render: function(){
		var page = this.props.page,
			status = this.props.status;

		var titlebar = status.message || convert.URLToReadable(page.url);

		var flags = null;
		if(page.flags && page.flags.length > 0){
			flags = React.DOM.div({
				className: 'page-flags'
			}, page.flags.map(function(flag){
				return React.DOM.span({key: flag}, flag);
			}));
		}

		var note = null;
		if(status.text &&  !this.state.noteHidden){
			note = React.DOM.div({
				className: 'page-note',
				onClick: this.hideNote
			}, status.text);
		}

		var wideglyph = 'fullscreen-enter';
		if(this.props.wide){
			wideglyph = 'fullscreen-exit';
		}

		return React.DOM.div({
			className: 'page-status page-status-' + status.state + ' ' + this.props.className
		},
			React.DOM.span({
				className: 'page-titlebar',
				title: titlebar
			}, titlebar),
			note,
			flags,
			newIcon(wideglyph, {
				title: 'Toggle Page Width',
				onClick: this.props.onToggleWidth
			}),
			newIcon('x', {
				className: 'icon-red',
				title: 'Close Page',
				onClick: this.props.onClose
			})
		);
	}
});
