package('kb.Lineup', function(exports) {
	'use strict';

	depends('Lineup.js');
	depends('Stage.View.js');

	exports.View = React.createClass({
		displayName: 'Lineup',

		getInitialState: function() {
			return {
				width: window.clientWidth
			};
		},

		render: function() {
			var self = this;
			var containerWidth = this.state.width;

			// try to calculate the best sizes for normal and wide stages
			var normal = containerWidth;
			var wide = containerWidth;
			if (containerWidth > 465 * 3.5) {
				normal = containerWidth * 0.25;
				wide = containerWidth * 0.50;
			} else if (containerWidth > 465 * 2.5) {
				normal = containerWidth * 0.33;
				wide = containerWidth * 0.663;
			} else if (containerWidth > 465 * 1.5) {
				normal = containerWidth * 0.50;
				wide = containerWidth * 1.00;
			}
			normal = Math.min(normal, 500);
			wide = Math.min(wide, 700);

			var left = 0;
			return React.DOM.div({
					className: 'lineup'
				},
				this.props.Lineup.stages.map(function(stage) {
					var width = stage.wide ? wide : normal;
					var r = React.createElement(kb.Stage.View, {
						style: {
							width: width + 'px',
							left: left + 'px'
						},
						onWidthChanged: self.onStageWidthChanged,
						key: stage.id,
						stage: stage
					});
					left += width;
					return r;
				}));
		},

		// bindings to Lineup
		changed: function() {
			this.forceUpdate();
		},

		onStageWidthChanged: function() {
			this.forceUpdate();
		},
		resized: function() {
			this.setState({
				width: this.getDOMNode().clientWidth
			});
		},
		componentDidMount: function() {
			this.props.Lineup.on('changed', this.changed, this);
			window.onresize = this.resized;

			this.setState({
				width: this.getDOMNode().clientWidth
			});
		},
		componentWillReceiveProps: function(nextprops) {
			if (this.props.Lineup !== nextprops.Lineup) {
				this.props.Lineup.remove(this);
				nextprops.Lineup.on('changed', this.changed, this);
			}
		},
		componentWillUnmount: function() {
			window.onresize = null;
			this.props.Lineup.remove(this);
		}
	});
});
