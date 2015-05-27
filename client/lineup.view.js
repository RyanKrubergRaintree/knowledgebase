import "lineup.js";
import "stage.view.js";

KB.Lineup.View = (function(){
	var Lineup = React.createClass({
		displayName: "Lineup",

		render: function(){
			return React.DOM.div(
				{ className: "lineup" },
				this.props.Lineup.stages.map(function(stage){
					return React.createElement(KB.Stage.View, {
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
			this.props.Lineup.on("changed", this.changed, this);
		},
		componentWillReceiveProps: function(nextprops){
			if(this.props.Lineup !== nextprops.Lineup){
				this.props.Lineup.remove(this);
				nextprops.Lineup.on("changed", this.changed, this);
			}
		},
		componentWillUnmount: function() {
			this.props.Lineup.remove(this);
		}
	});

	return Lineup;
})();
