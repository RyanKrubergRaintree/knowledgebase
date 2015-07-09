this.Sanitize = (function(){
	"use strict";

	var iframe = document.createElement('iframe');
	if (iframe['sandbox'] === undefined) {
		return function(input){
			return input
		};
	}

	iframe.sandbox = 'allow-same-origin';
	iframe.style.display = 'none';
	document.body.appendChild(iframe);

	function clone(node){
		return node.cloneNode(true);
	}

	var body = iframe.contentDocument.body;
	return function(input){
		try {
			body.innerHTML = input;
			var result = clone(body);
			body.innerHTML = "";
			return result.innerHTML;
		} catch (ex) {
			return input;
		}
	};
})();
