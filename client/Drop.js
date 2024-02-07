package("kb.drop", function (exports) {
	"use strict";

	depends("Convert.js");
	var rxCode = /[=><;{}[\]]/;

	exports.Item = undefined;
	exports.Allowed = "";
	exports.Effect = "";

	exports.SetAllowed = function (ev, allowed) {
		exports.Allowed = allowed;
		try {
			ev.dataTransfer.allowedEffect = allowed;
		} catch (e) {
			/* empty */
		}
	};

	exports.GetAllowed = function (ev) {
		var result = "copy";
		try {
			result = ev.dataTransfer.allowedEffect;
		} catch (e) {
			result = exports.Allowed;
		}
		return result;
	};

	exports.SetEffect = function (ev, effect) {
		exports.Effect = effect;
		try {
			// TODO: understand if this is supposed to be undefined or something else
			// eslint-disable-next-line no-undef
			ev.dataTransfer.dropEffect = dropEffect;
		} catch (e) {
			/* empty */
		}
	};

	exports.GetEffect = function (ev) {
		var result = "copy";
		try {
			result = ev.dataTransfer.dropEffect;
		} catch (e) {
			result = exports.Effect;
		}
		return result;
	};

	exports.EffectFor = function (ev) {
		var effect = exports.GetAllowed(ev);
		if (effect === "copy") {
			return "copy";
		}
		if (ev.shiftKey) {
			return "copy";
		}
		return "move";
	};

	exports.SetItem = function (ev, item) {
		exports.Item = JSON.stringify(item);
		try {
			ev.dataTransfer.setData("Text", JSON.stringify(item));
		} catch (e) {
			/* empty */
		}
	};

	exports.GetItem = function (ev) {
		var result = null;
		try {
			var data = ev.dataTransfer.getData("Text");
			result = JSON.parse(data);
		} catch (e) {
			result = exports.Item;
		}
		exports.Item = null;
		return result;
	};

	exports.SetDragImage = function (ev, node, x, y) {
		try {
			if (ev.dataTransfer.setDragImage) {
				ev.dataTransfer.setDragImage(node, x, y);
			}
		} catch (e) {
			/* empty */
		}
	};

	function getImage(dataTransfer) {
		var acceptedImages = {
			"image/png": true,
			"image/jpeg": true
		};
		if (typeof dataTransfer.files === "undefined") {
			return null;
		}
		for (var i = 0; i < dataTransfer.files.length; i += 1) {
			var file = dataTransfer.files[i];
			if (acceptedImages[file.type]) {
				return file;
			}
		}
		return null;
	}

	function resizeImage(src, onResized) {
		var MaxWidth = 1024,
			MaxHeight = 1024;

		var image = new Image();
		image.onload = function () {
			var canvas = document.createElement("canvas");
			if (image.height > MaxHeight) {
				image.width *= MaxHeight / image.height;
				image.height = MaxHeight;
			}
			if (image.width > MaxWidth) {
				image.height *= MaxWidth / image.width;
				image.width = MaxWidth;
			}

			canvas.width = image.width;
			canvas.height = image.height;

			var ctx = canvas.getContext("2d");
			ctx.clearRect(0, 0, canvas.width, canvas.height);
			ctx.drawImage(image, 0, 0, image.width, image.height);

			onResized(canvas.toDataURL());
		};
		image.src = src;
	}

	exports.ConvertUnknown = ConvertUnknown;
	function ConvertUnknown(stage, after, dataTransfer) {
		var item = createItem(stage, dataTransfer);
		if (item) {
			stage.patch({
				type: "add",
				after: after,
				id: item.id,
				item: item
			});
		}
	}

	function createItem(stage, dataTransfer) {
		var image = getImage(dataTransfer);
		if (image) {
			var item = {
				id: GenerateID(),
				type: "image",
				text: "",
				url: ""
			};

			// do delayed loading
			var reader = new FileReader();
			reader.onload = function (ev) {
				resizeImage(ev.target.result, function (data) {
					stage.patch({
						type: "edit",
						id: item.id,
						item: {
							id: item.id,
							type: "image",
							text: "",
							url: data
						}
					});
				});
			};
			reader.readAsDataURL(image);
			return item;
		}

		try {
			var html = dataTransfer.getData("text/html");
			var href = dataTransfer.getData("text/uri-list");

			if (href) {
				if (html) {
					var rxTags = /<[^>]+>/g;
					html = html.replace(rxTags, "");
				} else {
					html = kb.convert.URLToReadable(href);
				}

				return {
					id: GenerateID(),
					type: "reference",
					title: html,
					url: href
				};
			}
		} catch (ex) {
			// this getData may fail in IE
		}
		try {
			var text = dataTransfer.getData("text/plain");
			if (text === "") {
				text = dataTransfer.getData("Text");
			}
		} catch (ex) {
			text = dataTransfer.getData("Text");
		}
		if (text) {
			if (text.match(rxCode)) {
				return {
					id: GenerateID(),
					type: "code",
					text: text
				};
			}

			return {
				id: GenerateID(),
				type: "paragraph",
				text: text
			};
		}

		console.log("Unhandled drop item:", JSON.parse(JSON.stringify(dataTransfer)));
	}
});
