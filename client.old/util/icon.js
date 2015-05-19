import React from 'react';

export function newIcon(glyph, props){
	props = props || {};
	props['className'] = props['className'] || '';
	props['className'] += ' oi icon'
	if(props['onClick'] || props['isButton'] ){
		props['className'] += ' icon-clickable';
	}
	if(props['draggable']) {
		props['className'] += ' icon-draggable';
	}
	props['data-glyph'] = glyph;
	return React.DOM.a(props, ' ', props['text']);
}
