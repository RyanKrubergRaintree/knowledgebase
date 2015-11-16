package('kb', function(exports) {
	'use strict';

	depends('util/Notifier.js');
	depends('Convert.js');
	depends('Stage.js');

	exports.Lineup = Lineup;

	function Lineup() {
		this.stages = [];
		this.notifier = new kb.util.Notifier();
		this.notifier.mixto(this);
	}

	Lineup.prototype = {
		changed: function() {
			this.notifier.emit({
				type: 'changed',
				lineup: this
			});
		},

		stageById: function(id) {
			for (var i = 0; i < this.stages.length; i += 1) {
				if (this.stages[i].id === id) {
					return this.stages[i];
				}
			}
			return undefined;
		},
		indexOf_: function(id) {
			if (typeof id === 'undefined') {
				return -1;
			}
			for (var i = 0; i < this.stages.length; i += 1) {
				if (this.stages[i].id === id) {
					return i;
				}
			}
			return -1;
		},

		trim_: function(id) {
			if (typeof id === 'undefined') {
				return;
			}
			var i = this.indexOf_(id);
			if (i >= 0) {
				this.stages = this.stages.slice(0, i + 1);
			}
		},

		clear: function() {
			this.removeListeners();
			this.stages = [];
			this.changed();
		},

		closeLast: function() {
			// always keep one stage open
			if (this.stages.length > 1) {
				this.stages[this.stages.length - 1].close();
			}
		},

		changeRef: function(id, stage) {
			var i = this.indexOf_(id);
			if (i >= 0) {
				var ref = this.stages[i];
				ref.url = kb.convert.URLToReadable(stage.url);
				ref.link = kb.convert.URLToLink(stage.link);
				ref.title = stage.title;
				this.changed();
			}
		},

		// url
		// title
		// link
		// after, optional
		// insteadOf, optional
		open: function(props) {
			this.trim_(props.after);
			var stage = new kb.Stage(props);

			var i = this.indexOf_(props.insteadOf);
			if (i >= 0) {
				this.stages[i].remove(this);
				this.stages[i] = stage;
			} else {
				this.stages.push(stage);
			}

			stage.on('closed', this.handleClose, this);
			stage.on('urlChanged', this.handleURLChanged, this);
			this.changed();
			return stage.id;
		},

		openLink: function(link) {
			this.open(kb.convert.LinkToReference(link));
		},

		handleClose: function(ev) {
			this.stages = this.stages.filter(function(stage) {
				return stage !== ev.stage;
			});
			this.changed();
		},
		handleURLChanged: function( /*ev*/ ) {
			this.changed();
		},

		removeListeners: function() {
			this.stages.map(function(stage) {
				stage.remove(this);
			});
		},
		addListeners: function() {
			var self = this;
			this.stages.map(function(stage) {
				stage.on('closed', self.handleClose, self);
				stage.on('urlChanged', self.handleURLChanged, self);
			});
		},

		updateRefs: function(nextstages) {
			this.removeListeners();

			var stages = this.stages.slice();
			var changed = false;

			var newstages = nextstages.map(function(stage) {
				var prev = stages.shift();

				if (prev) {
					var plink = kb.convert.ReferenceToLink(prev);
					var slink = kb.convert.ReferenceToLink(stage);
					if (plink === slink) {
						return prev;
					}
				}
				changed = true;
				return stage;
			});

			if (stages.length > 0) {
				changed = true;
			}
			if (changed) {
				this.stages = newstages;
				this.changed();
			}

			this.addListeners();
		},

		findStageFromElement: function(el) {
			for (var i = 0; i < 64; i += 1) {
				if (el === null) {
					return null;
				}
				if (getClassList(el).contains('stage')) {
					var id = GetDataAttribute(el, 'id');
					return this.stageById(id);
				}
				el = el.parentElement;
			}
			return undefined;
		},

		handleOpenLink: function(ev) {
			ev.preventDefault();
			ev.stopPropagation();

			var target = ev.target;
			var stage = this.findStageFromElement(target);

			var ref = kb.convert.LinkToReference(target.href);
			var url = ref.url;

			if (stage) {
				var locFrom = kb.convert.URLToLocation(stage.url);
				var locTo = kb.convert.URLToLocation(url);
				if (locFrom.host === '') {
					url = locTo.path;
				} else {
					url = '//' + locFrom.host + locTo.path;
				}
			}

			var link = GetDataAttribute(target, 'link');
			var link = link || ref.link;
			var title = target.innerText;

			if (ev.ctrlKey) {
				this.open({
					url: url,
					link: link,
					title: title
				});
			} else {
				this.open({
					url: url,
					link: link,
					title: title,
					after: stage && stage.id
				});
			}
		},

		handleClickLink: function(ev) {
			var t = ev.target;
			if (t.nodeName !== 'A') {
				return;
			}
			if (getClassList(t).contains('external-link')) {
				return;
			}
			if ((t.onclick != null) || (t.onmousedown != null) || (t.onmouseup != null)) {
				return;
			}

			var href = t.attributes['href'];
			if (typeof href === 'undefined') {
				return;
			}
			var path = href.value;
			if ((path === '') || (path === '/') || (path === '#')) {
				return;
			}

			this.handleOpenLink(ev);
		}
	};
});
