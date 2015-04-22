'use strict';

if (!Object.clone) {
	Object.defineProperty(Object, 'clone', {
		enumerable: false,
		configurable: true,
		writable: true,
		value: function (obj) {
			var copy = JSON.parse(JSON.stringify(obj));
			copy.constructor = obj.constructor;
			copy.prototype = obj.prototype;
			copy.__proto__ = obj.__proto__;
			return copy;
		}
	});
}

if (!Object.assign) {
	Object.defineProperty(Object, 'assign', {
		enumerable: false,
		configurable: true,
		writable: true,
		value: function (target) {
			if (target === undefined || target === null) {
				throw new TypeError('Cannot convert first argument to object');
			}

			var to = Object(target);
			for (var i = 1; i < arguments.length; i++) {
				var nextSource = arguments[i];
				if (nextSource === undefined || nextSource === null) {
					continue;
				}
				var keysArray = Object.keys(Object(nextSource));
				for (var nextIndex = 0, len = keysArray.length; nextIndex < len; nextIndex++) {
					var nextKey = keysArray[nextIndex];
					var desc = Object.getOwnPropertyDescriptor(nextSource, nextKey);
					if (desc !== undefined && desc.enumerable) {
						to[nextKey] = nextSource[nextKey];
					}
				}
			}
			return to;
		}
	});
}

window.iff = function (v, e, o) {
	return v ? e : o || null;
};

// show the string that parse fails with
var _originalJSONParse = JSON.parse;
JSON.parse = function (data) {
	try {
		var result = _originalJSONParse.call(JSON, data);
	} catch (err) {
		console.error('JSON.parse failed: %s', data);
		throw err;
	}
	return result;
};
