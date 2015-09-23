package('kb.Page', function(exports) {
	'use strict';

	depends('Page.js');
	depends('Drop.js');
	depends('item/View.js');

	depends('util/ParseJSON.js');

	// ensure that we clear all the place-holders
	window.ondragend = clearDropPosition;
	window.ondrop = clearDropPosition;

	function clearDropPosition() {
		var els = document.querySelectorAll('.drop-after');
		for (var i = 0; i < els.length; i += 1) {
			getClassList(els[i]).remove('drop-after');
		}
		els = document.querySelectorAll('.drop-before');
		for (i = 0; i < els.length; i += 1) {
			getClassList(els[i]).remove('drop-before');
		}
	}

	function getAfter(page, drop) {
		if (drop === null) {
			return;
		}
		var story = page.story;

		if (drop.node === null) {
			if (drop.rel === 'after') {
				if (story.length > 0) {
					return story[story.length - 1].id; // as last
				}
			}
			return; // as first
		}

		var id = GetDataAttribute(drop.node, 'id');
		if (id === null) {
			return; // as first
		}
		for (var i = 0; i < story.length; i += 1) {
			if (story[i].id === id) {
				if (drop.rel === 'after') {
					return id;
				} else if (i > 0) {
					// as first
					return story[i - 1].id;
				} else {
					return;
				}
			}
		}
		return;
	}

	function findListContainer(node, containerClass, listClass) {
		while (node) {
			if (getClassList(node).contains(listClass)) {
				return node;
			}
			if (getClassList(node).contains(containerClass)) {
				return node.querySelector('.' + listClass);
			}
			node = node.parentElement;
		}
		return null;
	}

	function findListPosition(ev, containerClass, listClass, skip) {
		var list = findListContainer(ev.target, containerClass, listClass);
		if (list === null) {
			return;
		}

		var mouse = {
			x: ev.pageX,
			y: ev.pageY
		};


		var child = list.firstChild;
		if (child) {
			var box = child.getBoundingClientRect();
			if (child && (mouse.y < box.top)) {
				return {
					node: null,
					rel: 'before'
				};
			}
		} else {
			return {
				node: null,
				rel: 'before'
			};
		}

		for (var i = 0; i < list.children.length; i += 1) {
			var child = list.children[i];
			if (skip && skip(child)) {
				continue;
			}

			var box = child.getBoundingClientRect();
			if (mouse.y > box.bottom) {
				continue;
			}
			return {
				node: child,
				rel: 'after'
			};
		}

		return {
			node: null,
			rel: 'after'
		};
	}

	exports.View = React.createClass({
		displayName: 'Page',

		createReference: function(ev) {
			var stage = this.props.stage,
				page = stage.page;
			var item = {
				id: GenerateID(),
				type: 'reference',
				url: stage.url,
				title: page.title,
				text: page.synopsis
			};

			ev.dataTransfer.effectAllowed = 'copy';
			var data = {
				item: item
			};

			ev.dataTransfer.setData('Text', JSON.stringify(data));
		},

		dropEffectFor: function(ev) {
			var effect = 'copy';
			try {
				effect = ev.dataTransfer.effectAllowed;
			} catch (ex) {
				// HACK-FIX, this is required for IE11
				// otherwise getting effectAllowed fails
			}
			if (effect === 'copy') {
				return 'copy';
			}
			if (ev.shiftKey) {
				return 'copy';
			}
			return 'move';
		},

		dragEnter: function( /*ev*/ ) {},
		dragOver: function(ev) {
			clearDropPosition();
			if (!this.props.stage.canModify()) {
				return;
			}

			var drop = findListPosition(ev, 'page', 'page-story');
			if (drop === null) {
				return;
			}

			ev.preventDefault();
			ev.dataTransfer.dropEffect = this.dropEffectFor(ev);

			var drop = findListPosition(ev, 'page', 'page-story');
			if (drop === null) {
				return;
			}

			if (drop.node !== null) {
				getClassList(drop.node).add('drop-' + drop.rel);
			} else {
				var story = this.refs.story.getDOMNode();
				getClassList(story).add('drop-' + drop.rel);
			}
		},
		dragDrop: function(ev) {
			ev.preventDefault();
			ev.stopPropagation();

			var stage = this.props.stage,
				page = stage.page;

			clearDropPosition();
			if (!stage.canModify()) {
				return;
			}

			var drop = findListPosition(ev, 'page', 'page-story');
			if (drop === null) {
				return;
			}
			var after = getAfter(page, drop);

			var dropEffect = this.dropEffectFor(ev);

			var data = ev.dataTransfer.getData('Text');
			try {
				JSON.parse(data);
			} catch (ex) {
				data = null;
			}
			if (data) {
				var data = kb.util.ParseJSON(data);
				var item = data.item;

				if (data.url && (data.url !== page.url)) {
					item.origin = data.url;
					item.originId = item.id;
				}

				// if we make a copy or move to another page
				// we should update the id in the process
				if ((dropEffect === 'copy') || (data.url !== stage.url)) {
					item.id = GenerateID();
				}

				// are we moving on the same page?
				if ((dropEffect === 'move') && (data.url === stage.url)) {
					kb.item.DropCanceled = true;
					stage.patch({
						type: 'move',
						id: item.id,
						after: after
					});

					ev.dataTransfer.effectAllowed = 'none';
					ev.dataTransfer.dropEffect = 'none';
					return;
				}

				stage.patch({
					type: 'add',
					id: item.id,
					item: item,
					after: after
				});
				ev.dataTransfer.dropEffect = dropEffect;
			} else {
				ev.dataTransfer.dropEffect = 'copy';
				kb.DropData(stage, after, ev.dataTransfer);
			}
		},
		dragLeave: function( /*ev*/ ) {
			//TODO: fix glitchy rendering
			clearDropPosition();
		},
		render: function() {
			var stage = this.props.stage,
				owner = kb.convert.LinkToOwner(stage.link || stage.slug || stage.url || ''),
				page = this.props.page;

			var status;
			if (stage.state === 'loading') {
				status = React.DOM.div({
					className: 'page-loading'
				});
			} else if (stage.state !== 'loaded') {
				status = React.DOM.div({
						className: 'page-error'
					},
					React.DOM.h2(null, stage.lastStatusText),
					React.DOM.p(null, stage.lastError)
				);
			}

			return React.DOM.div({
					className: 'page',

					onDragEnter: this.dragEnter,
					onDragOver: this.dragOver,
					onDrop: this.dragDrop,
					onDragLeave: this.dragLeave
				},
				React.DOM.h2({
					className: 'page-owner'
				}, owner),
				React.DOM.h1({
					className: 'page-title',
					title: 'Drag to create a page reference.',
					draggable: true,
					onDragStart: this.createReference,
					style: {
						cursor: 'move'
					}
				}, page.title),
				status,
				React.createElement(Story, {
					ref: 'story',

					stage: stage,
					page: page,
					story: page.story
				})
			);
		}
	});

	var Story = React.createClass({
		displayName: 'Story',
		render: function() {
			var stage = this.props.stage,
				story = this.props.story;

			return React.DOM.div({
					className: 'page-story'
				},
				story.map(function(item, i) {
					return React.createElement(kb.item.View, {
						key: i + '|' + item.id,
						stage: stage,
						item: item
					});
				})
			);
		}
	});
});
