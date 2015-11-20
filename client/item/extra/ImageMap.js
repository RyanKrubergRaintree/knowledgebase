package("kb.item.content", function(exports) {
	"use strict";

	depends("ImageMap.css");

	exports["image-map"] = React.createClass({
		displayName: "ImageMap",
		contextTypes: {
			CurrentSelection: React.PropTypes.object
		},
		areaHoverStart: function(ev) {
			var id = GetDataAttribute(ev.target, "focusid");
			this.context.CurrentSelection.highlight(id);
		},
		areaHoverEnd: function(ev) {
			var id = GetDataAttribute(ev.target, "focusid");
			this.context.CurrentSelection.unhighlight(id);
		},
		areaSelect: function(ev) {
			var id = GetDataAttribute(ev.target, "focusid");
			this.context.CurrentSelection.toggleSelect(id);
		},
		render: function() {
			var item = this.props.item;
			var stage = this.props.stage;
			var size = item.size;
			var self = this;

			var loc = kb.convert.URLToLocation(stage.link);
			return React.DOM.div({
					className: "item-image-map content-image-map"
				},
				React.DOM.div({
						className: "container"
					},
					React.DOM.image({
						src: item.image
					}),
					React.DOM.div({
							className: "overlay",
							style: {
								maxWidth: size.x + "px"
							}
						},
						item.areas.map(function(area, index) {
							return React.DOM.a({
								key: index,
								"data-focusid": area.id,
								className: "area",
								style: {
									left: (area.min.x * 100 / size.x) + "%",
									top: (area.min.y * 100 / size.y) + "%",
									width: ((area.max.x - area.min.x) * 100 / size.x) + "%",
									height: ((area.max.y - area.min.y) * 100 / size.y) + "%"
								},
								title: area.alt,

								href: loc.path + "#" + area.id,
								"data-link": loc.path + "#" + area.id,

								onMouseEnter: self.areaHoverStart,
								onMouseLeave: self.areaHoverEnd,
								onMouseDown: self.areaSelect
							});
						}))
				),
				React.DOM.p({
					className: "tip"
				}, "Click any highlighted region for details.")
			);
		}
	});
});
