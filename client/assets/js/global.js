'use strict';

window.GetDataAttribute = function GetDataAttribute(el, name) {
	if (typeof el.dataset !== 'undefined') {
		return el.dataset[name];
	} else {
		return el.getAttribute('data-' + name);
	}
};

window.GenerateID = function GenerateID() {
	return Math.random().toString(16).substr(2) +
		Math.random().toString(16).substr(2);
};

window.TestCase = function TestCase(casename, runcase) {
	var assert = {
		'true': function(ok, msg) {
			if (!ok) {
				throw new Error(msg);
			}
		},
		'fail': function(err) {
			throw new Error(err);
		},
		'equal': function(actual, expect, msg) {
			if (actual !== expect) {
				var full = '\ngot ' + actual + '\nexp ' + expect;
				if (typeof msg !== 'undefined') {
					full = msg + full;
				}
				throw new Error(full);
			}
		}
	};

	try {
		runcase(assert);
	} catch (err) {
		console.error('assert ' + casename + ' failed:', err);
	}
};

window.getClassList = function(el) {
	function split(s) {
		return s.length ? s.split(/\s+/g) : [];
	}

	if ('classList' in el) {
		return el.classList;
	}

	return {
		add: function(token) {
			el.className += ' ' + token;
		},
		remove: function(token) {
			var tokens = ' ' + el.className + ' ';
			tokens = tokens.replace(' ' + token + ' ', '');
			el.className = tokens.trim();
		},
		contains: function(token) {
			return split(el.className).indexOf(token) >= 0;
		}
	};
};
