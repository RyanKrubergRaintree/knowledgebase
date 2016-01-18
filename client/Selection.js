package("kb", function(exports) {
	"use strict";

	depends("util/Notifier.js");

	exports.Selection = Selection;

	function Selection() {
		this.highlighted = "";
		this.selected = "";

		this.notifier = new kb.util.Notifier();
		this.notifier.mixto(this);
	}

	Selection.prototype = {
		changed: function() {
			this.notifier.emit({
				type: "changed",
				highlighted: this.highlighted,
				selected: this.selected
			});
		},
		select: function(id) {
			this.selected = id;
			this.changed();
		},
		unselect: function(id) {
			if (this.selected === id) {
				this.selected = "";
				this.changed();
			}
		},
		toggleSelect: function(id) {
			if (this.selected !== id) {
				this.selected = id;
			} else {
				this.selected = "";
			}
			this.changed();
		},
		highlight: function(id) {
			if (this.highlighted === id) {
				return;
			}
			this.highlighted = id;
			this.changed();
		},
		unhighlight: function(id) {
			if (this.highlighted === id) {
				this.highlighted = "";
				this.changed();
			}
		},


		_targetId: function(ev) {
			if ((!ev) || (!ev.target)) {
				return;
			}
			if (ev.target.nodeName !== "A") {
				return
			}
			var id = GetDataAttribute(ev.target, "focusid");
			if (id == null) {
				return;
			}
			return id;
		},

		highlightTarget: function(ev) {
			var id = this._targetId(ev);
			if (id) {
				this.highlight(id);
			}
		}
	};
});
