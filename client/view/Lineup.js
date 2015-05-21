//import "/view/Stage.js"
//import "/view/View.js"

View.Lineup = (function(){
	var Lineup = React.createClass({
		displayName: "Lineup",

		getInitialState: function(){
			return {
				stages: Global.Lineup.stages,
			}
		},

		componentDidMount: function(){
			Global.Lineup.on("changed", this.changed, this);
		},
		componentWillUnmount: function() {
			Global.Lineup.off("changed", this.changed, this);
		},

		changed: function() {
			this.setState({stages: Global.Lineup.stages});
		},
		render: function(){
			return React.DOM.div(
				{ className: "lineup" },
				this.state.stages.map(function(proxy){
					return React.createElement(View.Stage, {
						key: proxy.key,
						proxy: proxy
					});
				}
			));
		}
	});

	return Lineup;
})();
