package('kb.item', function(exports){
	'use strict';

	var iframe = document.createElement('iframe');
	if (typeof iframe.sandbox === 'undefined') {
		exports.Sanitize = function(input){
			return input;
		};
		return;
	}

	iframe.sandbox = 'allow-same-origin';
	iframe.style.display = 'none';
	document.body.appendChild(iframe);

	function clone(node){
		return node.cloneNode(true);
	}

	var body = iframe.contentDocument.body;
	exports.Sanitize = Sanitize;

	function Sanitize(input){
		try {
			body.innerHTML = input;
			var result = clone(body);
			body.innerHTML = '';
			return result.innerHTML;
		} catch (ex) {
			return input;
		}
	}
});
