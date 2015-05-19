'use strict';
import React from 'react';

import {newIcon} from 'util/icon';
import ObjectId from 'ObjectId';
import {DragContext} from 'view/Dragging';

export var Factory = React.createClass({
	displayName: 'Factory',

	context: null,
	startDrag: function(ev){
		var proxy = this.props.proxy;
		var context = new DragContext(proxy);

		var item = {
			id: ObjectId(),
			type: this.props.type,
			text: ''
		};

		if(this.props.type === 'reference'){
			item.title = proxy.page.title;
			item.url = proxy.page.path;
			item.text = proxy.page.synopsis;
		}

		context.start(ev, item, null, true);
	},
	render: function(){
		return React.DOM.div({
			className: 'factory',
			draggable: true,
			onDragStart: this.startDrag,
			title: this.props.title
		}, this.props.type);
	}
})
