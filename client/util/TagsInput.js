import React from 'react';

export var TagsInput = React.createClass({
	getInitialState: function(){
		var tags = this.props.defaultValue;
		return {
			tagtext: tags.join(' ')
		};
	},
	getTags: function(){
		return this.parseTags(this.state.tagtext, false);
	},
	parseTags: function(tagtext, dedupeExcept){
		// sanitize
		tagtext = tagtext.replace(/[^ a-z0-9\-\.]/g, ' ');
		tagtext = tagtext.replace(/ +/g, ' ');

		// split
		var tags = tagtext.split(' ');

		// dedupe
		var offset = 0;
		if(dedupeExcept){
			offset = 1;
		}
		for(var i = tags.length-1-offset; i >= 0; i -= 1){
			var x = tags.indexOf(tags[i]);
			if((x >= 0) && (x < i)){
				tags.splice(i, 1);
			}
		}
		// remove empty tags
		if(!dedupeExcept){
			tags = tags.filter(function(tag){ return tag != ''; });
		}
		return tags;
	},
	update: function(){
		var tagtext = this.refs.input.getDOMNode().value;
		var tags = this.parseTags(tagtext, true);
		this.setState({tagtext: tags.join(' ')});
	},
	render: function(){
		return React.DOM.input({
			className: this.props.className,
			autoFocus: this.props.autoFocus,

			ref: 'input',
			value: this.state.tagtext,
			onChange: this.update
		})
	}
});
