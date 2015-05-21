//import "/util/SmoothScroll.js"
//import "/kb/Stage.js"
//import "/kb/Page.View.js"

KB.Stage.View = (function(){
	var StageButtons = React.createClass({
		displayName: "StageButtons",

		toggleWidth: function(){
			this.props.onToggleWidth();
		},
		close: function(){
			this.props.stage.close();
		},

		render: function(){
			var stage = this.props.stage;
			var a = React.DOM.a;
			return React.DOM.div(
				{className: "stage-buttons"},
				a({
					className:"mdi mdi-playlist-plus",
					href:"#",
					title:"Add an item."
				}),
				a({
					className:"mdi " + (this.props.isWide ? "mdi-arrow-collapse" : "mdi-arrow-expand"),
					title:"Toggle page width.",
					onClick: this.toggleWidth
				}),
				a({
					className:"mdi mdi-close",
					title:"Close page.",
					onClick: this.close
				})
			);
		}
	});

	var StageInfo = React.createClass({
		displayName: "StageInfo",
		render: function(){
			var table = React.DOM.table,
			 	tr = React.DOM.tr,
			 	td = React.DOM.td;

			return table({className:"stage-info"},
				tr(null, td(null, "Link"),        td(null, this.props.stage.link)),
				tr(null, td(null, "Create by"),   td(null, "Raintree Systems Help")),
				tr(null, td(null, "Shared with"), td(null, "Everyone"))
			);
		}
	});

	var Stage = React.createClass({
		displayName: "Stage",

		getInitialState: function(){
			return {
				wide: false
			};
		},

		toggleWidth: function(){
			this.setState({
				wide: !this.state.wide
			});
		},

		activate: function(ev){
			if(typeof ev == 'undefined'){
				SmoothScroll.to(this.getDOMNode());
			} else if (!ev.defaultPrevented){
				SmoothScroll.to(this.getDOMNode());
			}
		},

		render: function(){
			var wide = this.state.wide ? " stage-wide" : "";
			return React.DOM.div(
				{className: "stage" + wide, onClick: this.activate},
				React.createElement(StageButtons, {
					stage: this.props.stage,
					isWide: this.state.wide,
					onToggleWidth: this.toggleWidth
				}),
				React.DOM.div(
					{className:"stage-scroll round-scrollbar"},
					React.createElement(StageInfo, this.props),
					React.createElement(KB.Page.View, {
						stage: this.props.stage,
						page: this.props.stage.page
					})
				)
			);
		}
	});

	return Stage;
})();
