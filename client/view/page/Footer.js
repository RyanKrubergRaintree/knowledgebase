'use strict';

import React from 'react';

import {Factory} from 'view/item/Factory';
import {newIcon} from 'util/icon';

export var Footer = React.createClass({
	displayName: 'footer',
	render: function(){
		var proxy = this.props.proxy,
			page = proxy.page;

		return React.DOM.div({
			className:'page-footer ' + this.props.className
		},
			React.createElement(Factory, {
				proxy: proxy,
				type: 'reference',
				title: 'Reference to this page.'
			}),
			iff(!page.readonly, React.createElement(Factory, {
				proxy: proxy,
				type: 'paragraph',
				title: 'A paragraph of text.'
			})),
			iff(!page.readonly, React.createElement(Factory, {
				proxy: proxy,
				type: 'html',
				title: 'HTML content.'
			})),
			iff(!page.readonly, React.createElement(Factory, {
				proxy: proxy,
				type: 'code',
				title: 'Code block.'
			})),
			iff(!page.dynamic, newIcon('comment-square', {
				style: { 'float': 'right' },
				title: 'Show/Hide Comments',
				onClick: this.props.onToggleComments,
				text: page.comments.length
			}))
		);
	}
})
