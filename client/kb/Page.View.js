//import "/util/SmoothScroll.js"
//import "/kb/Page.js"
//import "/kb/ItemView.js"

KB.Page.View = (function(){
	var Page = React.createClass({
		displayName: "Page",
		render: function(){
			var stage = this.props.stage,
				page = this.props.page;

			return React.DOM.div(
				{className: "page"},
				React.DOM.h1(null, "Hello World"),
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
					return React.createElement(Item, {
						key: item.id || i,
						stage: stage,
						item: item
					});
				})
			);
		}
	});

	var Item = React.createClass({
		displayName: "Item",
		render: function(){
			var stage = this.props.stage,
				item = this.props.item;

			var view = ItemView[item.type] || ItemView.Unknown;
			return React.DOM.div(
				{ className: "item" },
				React.DOM.div({
					className:"item-drag",
					title: "Move or copy item."
				}),
				React.createElement(view, {
					stage: stage,
					item: item
				})
			);
		}
	});

	return Page;
})();
