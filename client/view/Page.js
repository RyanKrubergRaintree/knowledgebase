//import "/view/View.js"
//import "/util/SmoothScroll.js"

View.Page = (function(){
	var PageButtons = React.createClass({
		displayName: "PageButtons",
		render: function(){
			var a = React.DOM.a;
			return React.DOM.div(
				{className: "page-buttons"},
				a({className:"mdi mdi-playlist-plus", href:"#", title:"Add an item."}),
				a({className:"mdi mdi-arrow-expand", href:"#", title:"Toggle page width."}),
				a({className:"mdi mdi-close", href:"#", title:"Close page."})
			);
		}
	});

	var PageMeta = React.createClass({
		displayName: "PageMeta",
		render: function(){
			var table = React.DOM.table,
			 	tr = React.DOM.tr,
			 	td = React.DOM.td;

			return table({className:"page-meta"},
				tr(null, td(null, "Link"), td(null, this.props.proxy.link)),
				tr(null, td(null, "Create by"), td(null, "Raintree Systems Help")),
				tr(null, td(null, "Shared with"), td(null, "Everyone"))
			);
		}
	});

	var PageContent = React.createClass({
		displayName: "PageContent",
		render: function(){
			return React.DOM.div(
				{className: "page-content"},
				React.DOM.h1(null, "Hello World"),
				React.createElement(PageStory, {})
			);
		}
	});

	var PageStory = React.createClass({
		displayName: "PageStory",
		render: function(){
			return React.DOM.div(
				{className: "page-story"},
				React.DOM.p({}, "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Cum rem accusantium libero eligendi repellat, quae commodi debitis odit animi facere illum! Laboriosam fugiat iste accusamus vitae doloremque corporis, nisi dolor."),
				React.DOM.p({}, "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Cum rem accusantium libero eligendi repellat, quae commodi debitis odit animi facere illum! Laboriosam fugiat iste accusamus vitae doloremque corporis, nisi dolor."),
				React.DOM.p({}, "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Cum rem accusantium libero eligendi repellat, quae commodi debitis odit animi facere illum! Laboriosam fugiat iste accusamus vitae doloremque corporis, nisi dolor."),
				React.DOM.p({}, "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Cum rem accusantium libero eligendi repellat, quae commodi debitis odit animi facere illum! Laboriosam fugiat iste accusamus vitae doloremque corporis, nisi dolor.")
			);
		}
	});

	var Page = React.createClass({
		displayName: "Page",

		activate: function(ev){
			if(typeof ev == 'undefined'){
				SmoothScroll.to(this.getDOMNode());
			} else if (!ev.defaultPrevented){
				SmoothScroll.to(this.getDOMNode());
			}
		},

		render: function(){
			return React.DOM.div(
				{className: "page-size", onClick: this.activate},
				React.createElement(PageButtons, {}),
				React.DOM.div(
					{className:"page-scroll round-scrollbar"},
					React.createElement(PageMeta, {proxy: this.props.proxy}),
					React.createElement(PageContent, {})
				)
			);
		}
	});

	return Page;
})();
