import "util/SmoothScroll.js";
import "stage.js";
import "page.view.js";

KB.Stage.View = (function(){
	var StageButtons = React.createClass({
		displayName: "StageButtons",

		toggleWidth: function(){
			this.props.onToggleWidth();
		},
		close: function(){
			this.props.stage.close();
		},

		createFactory: function(ev){
			var item = {
				id: GenerateID(),
				type: "factory",
				text: ""
			};

			ev.dataTransfer.effectAllowed = 'copy';
			var data = {
				item: item
			};

			ev.dataTransfer.setData("kb/item", JSON.stringify(data));
		},

		render: function(){
			var stage = this.props.stage;
			var a = React.DOM.a;
			return React.DOM.div(
				{className: "stage-buttons"},
				a({
					className:"mdi mdi-playlist-plus",
					title:"Drag to page to add an item.",
					style: { cursor: "move" },
					draggable: true,
					onDragStart: this.createFactory
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
				tr(null, td(null, "Link"),  td(null, this.props.stage.link)),
				tr(null, td(null, "State"), td(null, this.props.stage.state))
			);
		}
	});

	function extractGroup(link) {
		if(link == null){
			return "";
		}
		var i = link.indexOf(":")
		if(i >= 0) {
			return link.substr(0, i);
		}
		return "";
	}

	var NewPage = React.createClass({
		displayName: "NewPage",
		tryCreate: function(ev){
			ev.preventDefault();
			ev.stopPropagation();
		},
		render: function(){
			var stage = this.props.stage;
			var groups = [
				{id: "community", name: "Community"},
				{id: "engineering", name: "Engineering"}
			];
			return React.DOM.div(
				{ className: "page new-page" },
				React.DOM.form({
					onSubmit: this.tryCreate
				},
					React.DOM.label({
						htmlFor: "new-page-title",
					}, "Title"),
					React.DOM.input({
						id: "new-page-title",
						className: "title",
						defaultValue: stage.title,
						autoFocus: true
					}),
					React.DOM.label({}, "Page owner"),
					React.DOM.div(
						{ className: "group" },
						groups.map(function(group, i){
							return (
								React.DOM.div(
									{ key: group.id },
									React.DOM.input({
										id: "group-" + group.id,
										type: 'radio',
										name: 'group',
										value: group.id,
										defaultChecked: i == 0
									}),
									React.DOM.label({
										htmlFor: "group-" + group.id
									}, group.name)
								)
							);
						})
					),
					React.DOM.button({ type: "submit" }, "Create")
				)
			);
		}
	})

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
			window.setTimeout(this.activate, 100);
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
			var stage = this.props.stage;

			var creating = stage.canCreate &&
				((stage.url == null) || (stage.state == "not-found"));

			if(creating){
				return React.DOM.div(
					{
						className: "stage" + wide,
						onClick: this.activate,
						'data-id': stage.id
					},
					React.createElement(StageButtons, {
						stage: this.props.stage,
						isWide: this.state.wide,
						onToggleWidth: this.toggleWidth
					}),
					React.DOM.div(
						{ className:"stage-scroll round-scrollbar"},
						React.createElement(StageInfo, this.props),
						React.createElement(NewPage, {
							stage: this.props.stage
						})
					)
				);
			}

			return React.DOM.div(
				{
					className: "stage" + wide,
					onClick: this.activate,
					'data-id': stage.id
				},
				React.createElement(StageButtons, {
					stage: this.props.stage,
					isWide: this.state.wide,
					onToggleWidth: this.toggleWidth
				}),
				React.DOM.div(
					{ className:"stage-scroll round-scrollbar"},
					React.createElement(StageInfo, this.props),
					React.createElement(KB.Page.View, {
						stage: this.props.stage,
						page: this.props.stage.page,
					})
				)
			);
		},

		// bindings to Stage
		changed: function(){
			this.forceUpdate();
		},
		componentDidMount: function(){
			this.props.stage.on("changed", this.changed, this);
			this.props.stage.pull();
			this.activate();
		},
		componentWillReceiveProps: function(nextprops){
			if(this.props.stage !== nextprops.stage){
				this.props.stage.remove(this);
				nextprops.stage.on("changed", this.changed, this);
				nextprops.stage.pull();
			}
		},
		componentWillUnmount: function() {
			this.props.stage.remove(this);
		}
	});

	return Stage;
})();
