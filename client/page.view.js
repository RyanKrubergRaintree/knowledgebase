import "util/SmoothScroll.js";
import "page.js";
import "item.view.js";

KB.Page.View = (function(){
	var Page = React.createClass({
		displayName: "Page",

		render: function(){
			var stage = this.props.stage,
				page = this.props.page;

			return React.DOM.div(
				{ className: "page" },
				React.DOM.h1(null, page.title),
				React.createElement(Story, {
					stage: stage,
					page: page,
					story: page.story
				})
			);
		}
	});

	var Story = React.createClass({
		displayName: "Story",
		render: function(){
			var stage = this.props.stage,
				page = this.props.page,
				story = this.props.story;

			return React.DOM.div(
				{className: "page-story"},
				story.map(function(item, i){
					return React.createElement(KB.Item.View, {
						key: item.id || i,
						stage: stage,
						item: item
					});
				})
			);
		}
	});

	return Page;
})();
