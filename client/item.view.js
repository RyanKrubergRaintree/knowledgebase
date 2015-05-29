import "convert.js";
import "resolve.js";
import "item.js";
import "item.editor.js";

'use strict';

KB.Item.Content = {};

window.DropCanceled = null;

KB.Item.View = React.createClass({
	displayName: "Item",

	dragStart: function(ev, node, item){
		DropCanceled = false;
		var stage = this.props.stage,
			item = this.props.item;

		if(stage.canModify){
			ev.dataTransfer.effectAllowed = 'all';
		} else {
			ev.dataTransfer.effectAllowed = 'copy';
		}

		var off = mouseOffset(ev);
		ev.dataTransfer.setDragImage(this.getDOMNode(), off.x, off.y);

		var data = {
			item: item,
			title: stage.page.title,
			url: stage.url,
			text: stage.page.synopsis
		};
		ev.dataTransfer.setData("kb/item", JSON.stringify(data));
		ev.dataTransfer.setData("kb/url", stage.url);

		function mouseOffset(ev){
			ev = ev.nativeEvent || ev;
			return {
				x: ev.offsetX || ev.layerX || 0,
				y: ev.offsetY || ev.layerY || 0
			};
		}
	},
	drag: function(ev){ ev.preventDefault(); },
	dragEnd: function(ev){
		if(window.DropCanceled){
			ev.preventDefault();
			return;
		}
		var stage = this.props.stage,
			item = this.props.item;

		if(ev.dataTransfer.dropEffect == 'move'){
			stage.patch({
				type: 'remove',
				id: item.id
			});
		}

		ev.preventDefault();
		ev.stopPropagation();
	},

	startEditing: function(ev){
		var stage = this.props.stage,
			item = this.props.item;

		stage.editing.start(item.id);

		ev.preventDefault();
		ev.stopPropagation();
	},

	render: function(){
		var stage = this.props.stage,
			item = this.props.item;

		var view = KB.Item.Content[item.type] || KB.Item.Content.Unknown;
		var editingClass = '';
		var isEditing = false;
		if(stage.editing.item(item.id)){
			view = KB.Item.Editor;
			editingClass = ' item-editing';
			isEditing = true;
		}

		return React.DOM.div(
			{
				className: "item" + editingClass,
				onDoubleClick: stage.canModify ? this.startEditing : null,
				"data-id": item.id
			},
			!isEditing ? React.DOM.div({
				className:"item-drag",
				title: "Move or copy item.",
				draggable: true,

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

KB.Item.Content.Unknown = React.createClass({
	displayName: 'Unknown',
	render: function(){
		var item = this.props.item;
		return React.DOM.div(
			{ className: 'item-content content-unknown' },
			React.DOM.span({style: {"float": "right"}}, item.type),
			React.DOM.p({}, item.text),
			React.DOM.div({className:"clear-fix"})
		);
	}
});

var ContentTypes = [
	{name: "Text", type: "paragraph", desc: "simple text paragraph"},
	{name: "Tags", type: "tags", desc: "tags for the page"},
	{name: "HTML", type: "html", desc: "a subset of html for more advanced content"},
	{name: "Code", type: "code", desc: "item especially designed for code"},
];

KB.Item.Content['factory'] = React.createClass({
	displayName: 'Factory',
	convert: function(ev){
		var type = ev.currentTarget.dataset.type;
		var stage = this.props.stage,
			item = this.props.item;

		stage.patch({
			type: "edit",
			id: item.id,
			item: {
				type: type,
				id: item.id,
				text: item.text
			}
		});

		stage.editing.start(item.id);
	},

	render: function(){
		var self = this;
		var item = this.props.item;
		return React.DOM.div(
			{ className: 'item-content content-factory'	},
			React.DOM.p({}, item.text || "Add "),
			ContentTypes.map(function(item){
				return React.DOM.button(
					{
						key: item.type,
						className: "factory-item",
						"data-type": item.type,
						title: item.desc,
						onClick: self.convert
					}, item.name);
			})
		);
	}
});

KB.Item.Content['image'] = React.createClass({
	displayName: 'Image',
	render: function(){
		return React.DOM.div({
			className: 'item-content content-image'
		},
			React.DOM.img({src: this.props.item.url}),
			React.DOM.p({}, this.props.item.text)
		);
	}
});

KB.Item.Content['paragraph'] = React.createClass({
	displayName: 'Paragraph',
	render: function(){
		var stage = this.props.stage;
		return React.DOM.p({
			className: 'item-content content-paragraph',
			dangerouslySetInnerHTML: {
				__html: Resolve(stage, this.props.item.text)
			}
		});
	}
});

KB.Item.Content['html'] = React.createClass({
	displayName: 'HTML',
	render: function(){
		var stage = this.props.stage;
		return React.DOM.div({
			className: 'item-content content-html',
			dangerouslySetInnerHTML: {
				__html: ResolveHTML(stage, this.props.item.text)
			}
		});
	}
});

KB.Item.Content['code'] = React.createClass({
	displayName: 'Code',
	render: function(){
		return React.DOM.div({
			className: 'item-content content-code'
		}, this.props.item.text);
	}
});

KB.Item.Content['reference'] = React.createClass({
	displayName: 'Reference',
	render: function(){
		var item = this.props.item;
		var url = item.url;
		var loc = Convert.URLToLocation(url);
		var external = loc.origin && (loc.origin != window.location.origin);

		return React.DOM.div({className: 'item-content content-reference'},
			React.DOM.a({
				className: external ? 'external-link': '',
				target: external ? '_blank': '',
				href: url
			}, item.title),
			React.DOM.p({}, this.props.item.text)
		);
	}
});

KB.Item.Content['entry'] = React.createClass({
	displayName: 'Entry',
	render: function(){
		var item = this.props.item;
		var ref = Convert.LinkToReference(item.link);
		var url = ref.url;
		return React.DOM.div({
			className: 'item-content content-entry',
		},
			React.DOM.a({
				className: 'entry-title',
				title: url,
				href: url
			}, item.title),
			React.DOM.div({className: 'entry-owner'}, ref.owner),
			React.DOM.p({className: 'entry-synopsis'}, this.props.item.text)
		);
	}
});

KB.Item.Content['tags'] = React.createClass({
	displayName: 'Tags',
	render: function(){
		var item = this.props.item,
			stage = this.props.stage;

		var tags = item.text == "" ? [] : item.text.split(",");
		return React.DOM.div({className: 'item-contet content-tags'},
			tags.length > 0 ?
				tags.map(function(tag, i){
					tag = tag.trim();
					return React.DOM.a({
						className: "tag",
						key: i,
						href: '/tag:' + Slugify(tag)
					}, tag);
				})
			: React.DOM.p({}, "Double click here to add page tags."),
			React.DOM.div({className:"clear-fix"})
		);
	}
});
