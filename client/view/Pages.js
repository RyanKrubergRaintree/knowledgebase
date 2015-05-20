//import "/view/Page.js"
//import "/view/View.js"

View.Pages = (function(){
	var Pages = React.createClass({
		displayName: "Pages",

		getInitialState: function(){
			return {
				proxies: Global.Lineup.proxies,
			}
		},

		componentDidMount: function(){
			Global.Lineup.on("changed", this.changed, this);
		},
		componentWillUnmount: function() {
			Global.Lineup.off("changed", this.changed, this);
		},

		changed: function() {
			this.setState({proxies: Global.Lineup.proxies});
		},
		render: function(){
			return React.DOM.div(
				{ className: "pages" },
				this.state.proxies.map(function(proxy){
					return React.createElement(View.Page, {
						key: proxy.key,
						proxy: proxy
					});
				}
			));
		}
	});

	return Pages;
})();
