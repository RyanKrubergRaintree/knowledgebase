package("kb.Lineup", function(exports) {
	"use strict";

	depends("Lineup.js");
	depends("Stage.View.js");
	depends("Pin.View.js");

	var SelectionStyle = createReactClass({
		displayName: "SelectionStyle",
		getInitialState: function() {
			return {
				selected: this.props.CurrentSelection.selected,
				highlighted: this.props.CurrentSelection.highlighted
			};
		},
		update: function(ev) {
			this.setState({
				selected: ev.selected,
				highlighted: ev.highlighted
			});
		},
		componentDidMount: function() {
			this.props.CurrentSelection.on("changed", this.update, this);
		},
		componentWillUnmount: function() {
			this.props.CurrentSelection.remove(this);
		},

		render: function() {
			var state = this.state;

			var style = "";
			if (state.highlighted !== "") {
				var id = state.highlighted;
				style += "[data-id=\"" + id + "\"], " +
					"[data-focusid=\"" + id + "\"]" +
					"{" +
					"outline: 1px dashed #22F !important;" +
					"background: rgba(127,127,255,0.1) !important;" +
					"}";
			}
			if (state.selected !== "") {
				var id = state.selected;
				style += "[data-id=\"" + id + "\"], " +
					"[data-focusid=\"" + id + "\"]" +
					"{" +
					"outline: 1px dashed #22F !important;" +
					"background: rgba(127,127,255,0.1) !important;" +
					"}";
			}

			return React.DOM.style({
				key: Math.random(),
				dangerouslySetInnerHTML: {
					__html: style
				}
			});
		}
	});

	exports.View = createReactClass({
		displayName: "Lineup",

		getInitialState: function() {
			return {
				width: window.clientWidth
			};
		},

		render: function() {
			var self = this;
			var containerWidth = this.state.width;

			var idealWidth = 550;

			var normal = idealWidth;
			var wide = 800;

			var fittedCount = containerWidth / idealWidth;
			if (fittedCount < 1) {
				normal = containerWidth;
			} else if (fittedCount < 3) {
				var targetCount = (fittedCount + 0.5) | 0 + 0.2;
				normal = Math.min(containerWidth / targetCount, idealWidth);
			}

			var left = 0;
			var stages = this.props.Lineup.stages.map(function(stage) {
				var width = stage.wide ? wide : normal;
				var r = React.createElement(kb.Stage.View, {
					style: {
						width: width + "px",
						left: left + "px"
					},
					onWidthChanged: self.onStageWidthChanged,
					key: stage.id,
					stage: stage
				});
				left += width;
				return r;
			});

			var pinned = null;
			if (this.props.Lineup.pinned.visible) {
				var pin = this.props.Lineup.pinned;
				pinned = React.createElement(kb.Pin.View, {
					style: {
						left: left + "px",
						width: pin.width + "px"
					},
					width: pin.width,
					height: pin.height,
					onHide: this.hidePin,
					url: pin.url
				});
			}

			return React.DOM.div({
					className: "lineup"
				},
				React.createElement(SelectionStyle, this.props),
				stages,
				pinned
			);
		},

		// bindings to Lineup
		changed: function() {
			this.forceUpdate();
		},

		hidePin: function() {
			this.props.Lineup.hidePin();
		},
		onStageWidthChanged: function() {
			this.forceUpdate();
		},
		resized: function() {
			this.setState({
				width: ReactDOM.findDOMNode(this).clientWidth
			});
		},
		componentDidMount: function() {
			this.props.Lineup.on("changed", this.changed, this);
			window.onresize = this.resized;

			this.setState({
				width: ReactDOM.findDOMNode(this).clientWidth
			});

		},
		componentWillReceiveProps: function(nextprops) {
			if (this.props.Lineup !== nextprops.Lineup) {
				this.props.Lineup.remove(this);
				nextprops.Lineup.on("changed", this.changed, this);
			}
		},
		componentWillUnmount: function() {
			window.onresize = null;
			this.props.Lineup.remove(this);
		}
	});
});
