package('kb.item.content', function(exports) {
	'use strict';

	depends('ImageMap.css');

	exports['image-map'] = React.createClass({
		displayName: 'ImageMap',
		areaHoverStart: function(ev) {
			var id = GetDataAttribute(ev.target, 'focusid');
			kb.app.CurrentSelection.highlight(id);
		},
		areaHoverEnd: function(ev) {
			var id = GetDataAttribute(ev.target, 'focusid');
			kb.app.CurrentSelection.unhighlight(id);
		},
		areaSelect: function(ev) {
			var id = GetDataAttribute(ev.target, 'focusid');
			kb.app.CurrentSelection.toggleSelect(id);
		},
		render: function() {
			var item = this.props.item;
			var stage = this.props.stage;
			var size = item.size;
			var self = this;

			var loc = kb.convert.URLToLocation(stage.link);
			return React.DOM.div({
					className: 'item-image-map content-image-map'
				},
				React.DOM.div({
						className: 'content-image-map image',
						style: {
							width: size.x,
							height: size.y,
							background: 'url("' + item.image + '")'
						}
					},
					item.areas.map(function(area, index) {
						return React.DOM.a({
							key: index,
							'data-focusid': area.id,
							className: 'area',
							style: {
								left: area.min.x + 'px',
								top: area.min.y + 'px',
								width: area.max.x - area.min.x + 'px',
								height: area.max.y - area.min.y + 'px'
							},
							title: area.alt,

							href: loc.path + '#' + area.id,
							'data-link': loc.path + '#' + area.id,

							onMouseEnter: self.areaHoverStart,
							onMouseLeave: self.areaHoverEnd,
							onMouseDown: self.areaSelect
						});
					})
				)
			);
		}
	});
});
