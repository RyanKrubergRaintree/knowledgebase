'use strict';

import React from 'react';
import {Page} from './Page';

export var Pages = React.createClass({
	displayName: 'Pages',

	getInitialState: function(){
		return {
			pagerefs: this.props.Lineup.pagerefs
		};
	},
	componentDidMount: function(){
		this.props.Lineup.listen(this.changed);
	},
	componentWillUnmount: function() {
		this.props.Lineup.unlisten(this.changed);
	},
	changed: function(pagerefs){
		this.setState({pagerefs: pagerefs});
	},
	render: function(){
		var self = this;
		return React.DOM.div({id: 'pages'},
			this.state.pagerefs.map(function(pageref){
				return React.createElement(Page, {
					key: pageref.key,
					pageref: pageref,
					Store: self.props.Store,
					Lineup: self.props.Lineup
				});
			}));
	}
});
