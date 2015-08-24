package('kb.Lineup', function(exports){
	'use strict';

	depends('Lineup.js');
	depends('Stage.View.js');

	exports.View = React.createClass({
		displayName: 'Lineup',

		render: function(){
			return React.DOM.div(
				{ className: 'lineup' },
				this.props.Lineup.stages.map(function(stage){
					return React.createElement(kb.Stage.View, {
						key: stage.id,
						stage: stage
					});
				}
			));
		},

		// bindings to Lineup
		changed: function() {
			this.forceUpdate();
		},
		componentDidMount: function(){
			this.props.Lineup.on('changed', this.changed, this);
		},
		componentWillReceiveProps: function(nextprops){
			if(this.props.Lineup !== nextprops.Lineup){
				this.props.Lineup.remove(this);
				nextprops.Lineup.on('changed', this.changed, this);
			}
		},
		componentWillUnmount: function() {
			this.props.Lineup.remove(this);
		}
	});
});