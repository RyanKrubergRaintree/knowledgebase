//import "/kb/Stage.View.js"

KB.Lineup.View = (function(){
	var Lineup = React.createClass({
		displayName: "Lineup",

		componentDidMount: function(){
			this.props.Lineup.on("changed", this.changed, this);
		},
		componentWillReceiveProps: function(nextprops){
			this.props.Lineup.remove(this);
			nextprops.Lineup.on("changed", this.changed, this);
		},
		componentWillUnmount: function() {
			this.props.Lineup.remove(this);
		},

		changed: function() {
			this.forceUpdate();
		},
		render: function(){
			return React.DOM.div(
				{ className: "lineup" },
				this.props.Lineup.stages.map(function(stage){
					return React.createElement(KB.Stage.View, {
						key: stage.key,
						stage: stage
					});
				}
			));
		}
	});

	return Lineup;
})();
