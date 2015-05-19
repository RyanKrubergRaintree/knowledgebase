'use strict'

import React from 'react';

import {Item} from './Item';

export var Story = React.createClass({
	render: function(){
		var proxy = this.props.proxy,
			shouldEdit = this.props.shouldEdit;

		var items = this.props.story.map(function(item, index){
			var editing = shouldEdit.indexOf(item.id) >= 0;
			return React.createElement(Item, {
				key: item.id || 'i' + index,
				editing: editing,

				item: item,
				proxy: proxy
			});
		});

		return React.DOM.div({
			className: 'story ' + this.props.className
		}, items);
	}
});
