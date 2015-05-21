//import "/util/SmoothScroll.js"
//import "/kb/Page.js"

KB.Page.View = (function(){
	var Page = React.createClass({
		displayName: "Page",
		render: function(){
			return React.DOM.div(
				{className: "page"},
				React.DOM.h1(null, "Hello World"),
				React.createElement(Story, {})
			);
		}
	});

	var Story = React.createClass({
		displayName: "Story",
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

	return Page;
})();
