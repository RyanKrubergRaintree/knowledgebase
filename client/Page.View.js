package("kb.Page", function(exports) {
	"use strict";

	depends("Page.js");
	depends("Drop.js");
	depends("item/View.js");

	depends("util/ParseJSON.js");

	// ensure that we clear all the place-holders
	window.ondragend = clearDropPosition;
	window.ondrop = clearDropPosition;

	function clearDropPosition() {
		var els = document.querySelectorAll(".drop-after");
		for (var i = 0; i < els.length; i += 1) {
			getClassList(els[i]).remove("drop-after");
		}
		els = document.querySelectorAll(".drop-before");
		for (i = 0; i < els.length; i += 1) {
			getClassList(els[i]).remove("drop-before");
		}
	}

	function getAfter(page, drop) {
		if (drop === null) {
			return;
		}
		var story = page.story;

		if (drop.node === null) {
			if (drop.rel === "after") {
				if (story.length > 0) {
					return story[story.length - 1].id; // as last
				}
			}
			return; // as first
		}

		var id = GetDataAttribute(drop.node, "id");
		if (id === null) {
			return; // as first
		}
		for (var i = 0; i < story.length; i += 1) {
			if (story[i].id === id) {
				if (drop.rel === "after") {
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
				return node.querySelector("." + listClass);
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
					rel: "before"
				};
			}
		} else {
			return {
				node: null,
				rel: "before"
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
				rel: "after"
			};
		}

		return {
			node: null,
			rel: "after"
		};
	}

	exports.View = createReactClass({
		displayName: "Page",

		createReference: function(ev) {
			var stage = this.props.stage,
				page = stage.page;
			var item = {
				id: GenerateID(),
				type: "reference",
				url: stage.url,
				title: page.title,
				text: page.synopsis
			};

			kb.drop.SetAllowed(ev, "copy");
			kb.drop.SetItem(ev, {
				item: item
			});
		},

		dragEnter: function( /*ev*/ ) {},
		dragOver: function(ev) {
			clearDropPosition();
			if (!this.props.stage.canModify()) {
				return;
			}

			var drop = findListPosition(ev, "page", "page-story");
			if (drop === null) {
				return;
			}

			ev.preventDefault();
			kb.drop.SetEffect(ev, kb.drop.EffectFor(ev));

			var drop = findListPosition(ev, "page", "page-story");
			if (drop === null) {
				return;
			}

			if (drop.node !== null) {
				getClassList(drop.node).add("drop-" + drop.rel);
			} else {
				var story = this.refs.story;
				getClassList(story).add("drop-" + drop.rel);
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

			var drop = findListPosition(ev, "page", "page-story");
			if (drop === null) {
				return;
			}
			var after = getAfter(page, drop);
			var dropEffect = kb.drop.EffectFor(ev);
			var data = kb.drop.GetItem(ev);
			if (data) {
				var item = data.item;

				if (data.url && (data.url !== page.url)) {
					item.origin = data.url;
					item.originId = item.id;
				}

				// if we make a copy or move to another page
				// we should update the id in the process
				if ((dropEffect === "copy") || (data.url !== stage.url)) {
					item.id = GenerateID();
				}

				// are we moving on the same page?
				if ((dropEffect === "move") && (data.url === stage.url)) {
					kb.item.DropCanceled = true;
					stage.patch({
						type: "move",
						id: item.id,
						after: after
					});

					kb.drop.SetAllowed(ev, "none");
					kb.drop.SetEffect(ev, "none");
					return;
				}

				stage.patch({
					type: "add",
					id: item.id,
					item: item,
					after: after
				});
				kb.drop.SetEffect(ev, dropEffect);
			} else {
				kb.drop.SetEffect(ev, "copy");
				kb.drop.ConvertUnknown(stage, after, ev.dataTransfer);
			}
		},
		dragLeave: function( /*ev*/ ) {
			//TODO: fix glitchy rendering
			clearDropPosition();
		},
		render: function() {
			var stage = this.props.stage,
				owner = kb.convert.LinkToOwner(stage.link || stage.slug || stage.url || ""),
				page = this.props.page;

			var status;
			if (stage.state === "loading") {
				status = React.DOM.div({
					className: "page-loading"
				});
			} else if (stage.state !== "loaded") {
				status = React.DOM.div({
						className: "page-error"
					},
					React.DOM.h2(null, stage.lastStatusText),
					React.DOM.p(null, stage.lastError)
				);
			}

			return React.DOM.div({
					className: "page",

					onDragEnter: this.dragEnter,
					onDragOver: this.dragOver,
					onDrop: this.dragDrop,
					onDragLeave: this.dragLeave
				},
				React.DOM.h2({
					className: "page-owner"
				}, owner),
				React.DOM.h1({
					className: "page-title",
					title: "Drag to create a page reference.",
					draggable: true,
					onDragStart: this.createReference,
					style: {
						cursor: "move"
					}
				}, page.title),
				status,
				React.createElement(Story, {
					ref: "story",

					stage: stage,
					page: page,
					story: page.story
				})
			);
		}
	});

	var Story = createReactClass({
		displayName: "Story",
		render: function() {
			var stage = this.props.stage,
				story = this.props.story;

			return React.DOM.div({
					className: "page-story"
				},
				story.map(function(item, i) {
					return React.createElement(kb.item.View, {
						key: i + "|" + item.id,
						stage: stage,
						item: item
					});
				})
			);
		}
	});
});