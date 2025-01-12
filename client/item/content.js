package("kb.item.content", function (exports) {
	"use strict";

	depends("dita.css");
	depends("content.css");

	depends("../Convert.js");

	depends("Sanitize.js");
	depends("Resolve.js");

	window.ToggleCollapse = function ToggleCollapse(event) {
		var target = event.target;
		while (target) {
			var classes = getClassList(target);
			if (classes.contains("collapse")) {
				classes.toggle("collapsed");
				return;
			}
			target = target.parentNode;
		}
	};

	function InjectToggleHandler(htmlContent) {
		var x = htmlContent.replace(
			"class='collapse-icon",
			"onclick='ToggleCollapse(event)' class='collapse-icon"
		);
		return x.replace('class="collapse-icon', "onclick='ToggleCollapse(event)' class=\"collapse-icon");
	}

	exports.Unknown = createReactClass({
		displayName: "Unknown",
		render: function () {
			var item = this.props.item;
			return React.DOM.div(
				{
					className: "item-content content-unknown"
				},
				React.DOM.span(
					{
						style: {
							float: "right"
						}
					},
					item.type
				),
				React.DOM.p({}, item.text),
				React.DOM.div({
					className: "clear-fix"
				})
			);
		}
	});

	var ContentTypes = [
		{
			name: "Text",
			type: "paragraph",
			desc: "simple text paragraph"
		},
		{
			name: "HTML",
			type: "html",
			desc: "a subset of html for more advanced content"
		},
		{
			name: "Code",
			type: "code",
			desc: "item especially designed for code"
		},
		{
			name: "Tags",
			type: "tags",
			desc: "tags for the page"
		},
		{
			name: "Separator",
			type: "separator",
			desc: "line separator"
		}
	];

	exports.factory = createReactClass({
		displayName: "Factory",
		convert: function (ev) {
			var type = GetDataAttribute(ev.currentTarget, "type");
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

		render: function () {
			var self = this;
			var item = this.props.item;
			return React.DOM.div(
				{
					className: "item-content content-factory"
				},
				React.DOM.p({}, item.text || "Create new "),
				ContentTypes.map(function (item) {
					return React.DOM.button(
						{
							key: item.type,
							className: "factory-item",
							"data-type": item.type,
							title: item.desc,
							onClick: self.convert
						},
						item.name
					);
				})
			);
		}
	});

	exports.image = createReactClass({
		displayName: "Image",
		render: function () {
			return React.DOM.div(
				{
					className: "item-content content-image"
				},
				React.DOM.img({
					src: this.props.item.url
				}),
				React.DOM.p({}, this.props.item.text)
			);
		}
	});

	exports.paragraph = createReactClass({
		displayName: "Paragraph",
		render: function () {
			var stage = this.props.stage;
			var resolved = kb.item.Resolve(stage, this.props.item.text);
			var paragraphs = resolved.split("\n\n");
			if (paragraphs.length > 1) {
				return React.DOM.div(
					{
						className: "item-content content-paragraph"
					},
					paragraphs.map(function (p, i) {
						return React.DOM.p({
							key: i,
							dangerouslySetInnerHTML: {
								__html: InjectToggleHandler(kb.item.Sanitize(p))
							}
						});
					})
				);
			} else {
				return React.DOM.p({
					className: "item-content content-paragraph",
					dangerouslySetInnerHTML: {
						__html: InjectToggleHandler(kb.item.Sanitize(paragraphs[0]))
					}
				});
			}
		}
	});

	exports.html = createReactClass({
		displayName: "HTML",
		render: function () {
			var stage = this.props.stage;
			return React.DOM.div({
				className: "item-content content-html",
				dangerouslySetInnerHTML: {
					__html: InjectToggleHandler(
						kb.item.Sanitize(kb.item.ResolveHTML(stage, this.props.item.text))
					)
				}
			});
		}
	});

	exports.code = createReactClass({
		displayName: "Code",
		render: function () {
			return React.DOM.div(
				{
					className: "item-content content-code"
				},
				this.props.item.text
			);
		}
	});

	exports.reference = createReactClass({
		displayName: "Reference",
		render: function () {
			var item = this.props.item;
			var url = item.url;
			var loc = kb.convert.URLToLocation(url);
			var external = loc.host !== "" && loc.host !== window.location.host;

			return React.DOM.div(
				{
					className: "item-content content-reference"
				},
				React.DOM.a(
					{
						className: external ? "external-link" : "",
						target: external ? "_blank" : "",
						href: url
					},
					item.title
				),
				React.DOM.p({}, this.props.item.text)
			);
		}
	});

	exports.entry = createReactClass({
		displayName: "Entry",
		render: function () {
			var item = this.props.item;
			var ref = kb.convert.LinkToReference(item.link);
			var url = ref.url;
			return React.DOM.div(
				{
					className: "item-content content-entry"
				},
				React.DOM.a(
					{
						className: "entry-title",
						title: url,
						href: url
					},
					item.title
				),
				React.DOM.div(
					{
						className: "entry-owner"
					},
					ref.owner
				),
				React.DOM.p({
					className: "entry-synopsis",
					dangerouslySetInnerHTML: {
						__html: this.props.item.text
					}
				})
			);
		}
	});

	exports.tags = createReactClass({
		displayName: "Tags",
		render: function () {
			var item = this.props.item;

			var text = typeof item.text === "undefined" ? "" : item.text.trim();
			var tags = [];
			if (text !== "") {
				tags = text.split(",");
			}
			tags = tags
				.map(function (tag) {
					return tag.trim();
				})
				.filter(function (tag) {
					return tag !== "";
				});

			var hasValidTag = tags.length > 0;

			tags = tags.filter(function (tag) {
				// hide tags starting with id/...
				return !/^id\//.test(tag);
			});

			return React.DOM.div(
				{
					className: "item-contet content-tags"
				},
				hasValidTag
					? tags.map(function (tag, i) {
							tag = tag.trim();
							return React.DOM.a(
								{
									className: "tag",
									key: i,
									href: "/tag=pages/" + kb.convert.TextToSlug(tag)
								},
								tag
							);
						})
					: React.DOM.p({}, "Double click here to add page tags."),
				React.DOM.div({
					className: "clear-fix"
				})
			);
		}
	});

	exports.separator = createReactClass({
		displayName: "Entry",
		render: function () {
			var item = this.props.item;

			if (item.text === "") {
				return React.DOM.div(
					{
						className: "item-content content-separator"
					},
					React.DOM.hr(null)
				);
			}
			return React.DOM.div(
				{
					className: "item-content content-separator"
				},
				React.DOM.table(
					{
						style: {
							width: "100%"
						}
					},
					React.DOM.tbody(
						null,
						React.DOM.tr(
							null,
							React.DOM.td(null, React.DOM.hr(null)),
							React.DOM.td(
								{
									style: {
										width: "1px",
										padding: "0 10px",
										whiteSpace: "nowrap"
									}
								},
								this.props.item.text
							),
							React.DOM.td(null, React.DOM.hr(null))
						)
					)
				)
			);
		}
	});
});
