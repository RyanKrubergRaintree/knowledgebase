package('kb.item.content', function(exports) {
	'use strict';

	depends('ImageMap.css');

	exports['image-map'] = React.createClass({
		displayName: 'ImageMap',
		render: function() {
			var item = this.props.item;
			var size = item.size;

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
							className: 'area',
							style: {
								left: area.min.x + 'px',
								top: area.min.y + 'px',
								width: area.max.x - area.min.x + 'px',
								height: area.max.y - area.min.y + 'px'
							},
							title: area.alt
						});
					})
				)
			);
		}
	});
});
