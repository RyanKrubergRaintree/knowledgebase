//import "/util/SmoothScroll.js"
//import "/kb/Stage.js"
//import "/kb/Page.View.js"

KB.Stage.View = (function(){
	var StageButtons = React.createClass({
		displayName: "StageButtons",
		render: function(){
			var a = React.DOM.a;
			return React.DOM.div(
				{className: "stage-buttons"},
				a({className:"mdi mdi-playlist-plus", href:"#", title:"Add an item."}),
				a({className:"mdi mdi-arrow-expand", href:"#", title:"Toggle page width."}),
				a({className:"mdi mdi-close", href:"#", title:"Close page."})
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

		activate: function(ev){
			if(typeof ev == 'undefined'){
				SmoothScroll.to(this.getDOMNode());
			} else if (!ev.defaultPrevented){
				SmoothScroll.to(this.getDOMNode());
			}
		},

		render: function(){
			return React.DOM.div(
				{className: "stage", onClick: this.activate},
				React.createElement(StageButtons, {}),
				React.DOM.div(
					{className:"stage-scroll round-scrollbar"},
					React.createElement(StageInfo, {
						stage: this.props.stage
					}),
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
