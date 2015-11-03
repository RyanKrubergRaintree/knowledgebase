package('kb.item.content', function(exports) {
	'use strict';

	depends('ImageMap.css');

	exports['image-map'] = React.createClass({
		displayName: 'ImageMap',
		areaHoverStart: function(ev) {
			this.props.stage.selection.highlight(ev.target.dataset.ditaid);
		},
		areaHoverEnd: function(ev) {
			this.props.stage.selection.unhighlight(ev.target.dataset.ditaid);
		},
		areaSelect: function(ev) {
			this.props.stage.selection.toggleSelect(ev.target.dataset.ditaid);
		},
		render: function() {
			var item = this.props.item;
			var size = item.size;
			var self = this;

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
						return React.DOM.div({
							key: index,
							'data-ditaid': area.id,
							className: 'area',
							style: {
								left: area.min.x + 'px',
								top: area.min.y + 'px',
								width: area.max.x - area.min.x + 'px',
								height: area.max.y - area.min.y + 'px'
							},
							title: area.alt,

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
