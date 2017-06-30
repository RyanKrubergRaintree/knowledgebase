package("kb.item", function(exports) {
	"use strict";

	depends("Editor.js");

	depends("content.js");

	depends("extra/SimpleForm.js");
	depends("extra/Index.js");
	depends("extra/ImageMap.js");

	exports.DropCanceled = null;
	exports.View = createReactClass({
		displayName: "Item",

		dragStart: function(ev, node, item) {
			kb.item.DropCanceled = false;
			var stage = this.props.stage,
				item = this.props.item;

			if (stage.canModify()) {
				kb.drop.SetAllowed(ev, "all");
			} else {
				kb.drop.SetAllowed(ev, "copy");
			}

			var off = mouseOffset(ev);
			var viewnode = ReactDOM.findDOMNode(this);
			kb.drop.SetDragImage(ev, viewnode, off.x, off.y);

			var data = {
				item: item,
				title: stage.page.title,
				url: stage.url,
				text: stage.page.synopsis
			};
			kb.drop.SetItem(ev, data);

			function mouseOffset(ev) {
				ev = ev.nativeEvent || ev;
				return {
					x: ev.offsetX || ev.layerX || 0,
					y: ev.offsetY || ev.layerY || 0
				};
			}
		},
		drag: function() {},
		dragEnd: function(ev) {
			if (kb.item.DropCanceled) {
				ev.preventDefault();
				return;
			}
			var stage = this.props.stage,
				item = this.props.item;

			if (kb.drop.EffectFor(ev) === "move") {
				stage.patch({
					type: "remove",
					id: item.id
				});
			}

			ev.preventDefault();
			ev.stopPropagation();
		},

		startEditing: function(ev) {
			var stage = this.props.stage,
				item = this.props.item;

			stage.editing.start(item.id);

			ev.preventDefault();
			ev.stopPropagation();
		},

		render: function() {
			var stage = this.props.stage,
				item = this.props.item;

			var view = kb.item.content[item.type] || kb.item.content.Unknown;
			var editingClass = "";
			var isEditing = false;
			if (stage.editing.item(item.id)) {
				view = kb.item.Editor;
				editingClass = " item-editing";
				isEditing = true;
			}

			return React.DOM.div({
					className: "item" + editingClass,
					onDoubleClick: stage.canModify() ? this.startEditing : null,
					"data-id": item.id
				}, !isEditing ? React.DOM.a({
					className: "item-drag",
					draggable: "true",
					tabIndex: -1,

					href: "#",
					onClick: function(ev) {
						ev = ev || window.event;
						ev.preventDefault();
					},

					onDragStart: this.dragStart,
					onDrag: this.drag,
					onDragEnd: this.dragEnd
				}) : null,
				React.createElement(view, {
					stage: stage,
					item: item
				})
			);
		}
	});
});
