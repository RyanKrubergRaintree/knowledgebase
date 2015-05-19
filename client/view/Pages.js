//import "/view/Page.js"
//import "/view/View.js"

View.Pages = (function(){
	var Pages = React.createClass({
		displayName: "Pages",
		render: function(){
			return React.DOM.div({
				className: "pages"
			}, React.createElement(View.Page, {}))
		}
	});

	return Pages;
})();
