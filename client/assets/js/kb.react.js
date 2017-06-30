var kb = window.kb || kb || {};
kb.react = kb.react || {};

(function(exports) {
	"use strict";

	function checkPrimitive(expectedType, props, propName, componentName, location, propFullName) {
		if (props[propName] == null) {
			if (props[propName] === null) {
				return new Error('The ' + location + ' `' + propFullName + '` is marked as required ' + ('in `' + componentName + '`, but its value is `null`.'));
			}
			return new Error('The ' + location + ' `' + propFullName + '` is marked as required ' + ('in `' + componentName + '`, but its value is `undefined`.'));
		}
		var propType = typeof props[propName];
		if (propType != expectedType) {
			return new Error('Invalid ' + location + ' `' + propFullName + '` of type ' + ('`' + propType + '` supplied to `' + componentName + '`, expected ') + ('`' + expectedType + '`.'));
		}
		return null;
	}

	exports.object = function(props, propName, componentName, location, propFullName) {
		return checkPrimitive("object", props, propName, componentName, location, propFullName);
	};

	exports.func = function(props, propName, componentName, location, propFullName) {
		return checkPrimitive("function", props, propName, componentName, location, propFullName);
	};

	exports.string = function(props, propName, componentName, location, propFullName) {
		return checkPrimitive("string", props, propName, componentName, location, propFullName);
	};

	exports.number = function(props, propName, componentName, location, propFullName) {
		return checkPrimitive("number", props, propName, componentName, location, propFullName);
	};
})(kb.react);
