this.Sanitize = (function(){
	"use strict";

	var iframe = document.createElement('iframe');
	if (iframe['sandbox'] === undefined) {
		alert('Your browser does not support sandboxed iframes. Please upgrade to a modern browser.');
		return '';
	}

	iframe.sandbox = 'allow-same-origin';
	iframe.style.display = 'none';
	document.body.appendChild(iframe);

	function clone(node){
		return node.cloneNode(true);
	}

	var body = iframe.contentDocument.body;
	return function(input){
		body.innerHTML = input;
		var result = clone(body);
		body.innerHTML = "";
		return result.innerHTML;
	};
})();
