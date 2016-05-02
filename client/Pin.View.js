package("kb.Pin", function(exports) {
	"use strict";

	depends("util/SmoothScroll.js");

	var PinButtons = React.createClass({
		displayName: "PinButtons",

		close: function() {
			this.props.onHide();
		},

		render: function() {
			var stage = this.props.stage;
			var a = React.DOM.a;
			return React.DOM.div({
					className: "stage-buttons"
				},
				a({
					className: "mdi mdi-close",
					title: "Close pinned image.",
					onClick: this.close
				})
			);
		}
	});

	var View = React.createClass({
		displayName: "Pin",
		activate: function(ev) {
			if (typeof ev === "undefined") {
				var node = ReactDOM.findDOMNode(this);
				kb.util.SmoothScroll.to(node);
			} else if (!ev.defaultPrevented) {
				var node = ReactDOM.findDOMNode(this);
				kb.util.SmoothScroll.to(node);
			}
		},
		componentDidMount: function() {
			this.activate();
		},
		render: function() {
			var stage = this.props.stage,
				story = this.props.story;

			return React.DOM.div({
					className: "stage pin",
					onClick: this.activate,
					style: this.props.style,
				},
				React.createElement(PinButtons, {
					onHide: this.props.onHide
				}),
				React.DOM.div({
						className: "stage-scroll round-scrollbar"
					},
					React.DOM.div({
							className: "page page-full"
						},
						React.DOM.img({
							className: "pinned-image",
							src: this.props.url
						})
					)
				)
			);
		}
	});
	exports.View = View;
});
