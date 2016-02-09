package("kb.item.content", function(exports) {
	"use strict";

	depends("SimpleForm.css");

	exports["simple-form"] = React.createClass({
		displayName: "SimpleForm",
		contextTypes: {
			Session: React.PropTypes.object
		},
		getInitialState: function() {
			return {
				error: "",
				message: ""
			};
		},

		done: function(ev) {
			var xhr = ev.currentTarget;
			if (xhr.readyState !== 4) {
				return;
			}

			if (xhr.status !== 200) {
				this.setState({
					message: "",
					error: xhr.responseText
				});
				return;
			}
			var message = "done";
			if (xhr.responseText !== "") {
				message = xhr.responseText;
			}
			this.setState({
				message: message,
				error: ""
			});
			this.props.stage.refresh();
		},
		errored: function() {
			this.setState({
				message: "",
				error: "Unknown error."
			});
		},

		click: function(ev) {
			this.setState({
				error: "",
				message: "processing..."
			});

			var items = this.props.item.items || [];
			var body = {};
			var self = this;
			items.map(function(item) {
				var id = item.id || item.label;
				var node = self.refs[id];
				if (typeof node !== "undefined") {
					var value = node.value;
					body[id] = value;
				}
			});
			items["action"] = GetDataAttribute(ev.currentTarget, "action");

			this.context.Session.fetch({
				method: "POST",
				url: this.props.item.url,
				ondone: this.done,
				onerror: this.errored,
				body: items
			});
		},

		render: function() {
			var item = this.props.item,
				message = this.state.message,
				error = this.state.error;

			var self = this;

			var items = item.items || [];
			return React.DOM.form({
					className: "item-content content-simple-form",
					onSubmit: function(ev) {
						ev.preventDefault();
					}
				},
				item.text ? React.DOM.p({}, item.text) : null,
				message !== "" ? React.DOM.p({
					className: "message"
				}, message) : null,
				error !== "" ? React.DOM.p({
					className: "error"
				}, error) : null,
				items.map(function(item, i) {
					switch (item.type) {
						case "field":
							return React.DOM.input({
								key: i,
								ref: item.id || item.label,
								name: item.id || item.label,
								placeholder: item.label
							});
						case "button":
							return React.DOM.button({
								key: i,
								"data-action": item.action,
								onClick: self.click
							}, item.caption);
						case "option":
							return React.DOM.select({
								key: i,
								ref: item.id,
								name: item.id
							}, item.values.map(function(value, i) {
								return React.DOM.option({
									key: i,
									value: value
								}, value);
							}));
					}
				})
			);
		}
	});
});
