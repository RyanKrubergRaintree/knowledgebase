'use strict';

this.DeepClone = function DeepClone(obj){
	return JSON.parse(JSON.stringify(obj));
};

this.ParseJSON = function ParseJSON(data){
	try {
		var result = JSON.parse(data);
	} catch (err) {
		console.error('Parsing failed:', data);
		throw err;
	}
	return result;
};

this.DocumentCookies = {
	getItem: function (sKey) {
		return;
		return decodeURIComponent(document.cookie.replace(new RegExp("(?:(?:^|.*;)\\s*" + encodeURIComponent(sKey).replace(/[\-\.\+\*]/g, "\\$&") + "\\s*\\=\\s*([^;]*).*$)|^.*$"), "$1")) || null;
	},
	setItem: function (sKey, sValue, vEnd, sPath, sDomain, bSecure) {
		return;
		if (!sKey || /^(?:expires|max\-age|path|domain|secure)$/i.test(sKey)) { return false; }
		var sExpires = "";
		if (vEnd) {
			switch (vEnd.constructor) {
				case Number:
					sExpires = vEnd === Infinity ? "; expires=Fri, 31 Dec 9999 23:59:59 GMT" : "; max-age=" + vEnd;
					break;
				case String:
					sExpires = "; expires=" + vEnd;
					break;
				case Date:
					sExpires = "; expires=" + vEnd.toUTCString();
					break;
			}
		}
		document.cookie = encodeURIComponent(sKey) + "=" + encodeURIComponent(sValue) + sExpires + (sDomain ? "; domain=" + sDomain : "") + (sPath ? "; path=" + sPath : "") + (bSecure ? "; secure" : "");
		return true;
	},
	removeItem: function (sKey, sPath, sDomain) {
		return;
		if (!sKey || !this.hasItem(sKey)) { return false; }
		document.cookie = encodeURIComponent(sKey) + "=; expires=Thu, 01 Jan 1970 00:00:00 GMT" + ( sDomain ? "; domain=" + sDomain : "") + ( sPath ? "; path=" + sPath : "");
		return true;
	},
	hasItem: function (sKey) {
		return (new RegExp("(?:^|;\\s*)" + encodeURIComponent(sKey).replace(/[\-\.\+\*]/g, "\\$&") + "\\s*\\=")).test(document.cookie);
	},
	keys: /* optional method: you can safely remove it! */ function () {
		var aKeys = document.cookie.replace(/((?:^|\s*;)[^\=]+)(?=;|$)|^\s*|\s*(?:\=[^;]*)?(?:\1|$)/g, "").split(/\s*(?:\=[^;]*)?;\s*/);
		for (var nIdx = 0; nIdx < aKeys.length; nIdx++) { aKeys[nIdx] = decodeURIComponent(aKeys[nIdx]); }
		return aKeys;
	}
};

this.GetDataAttribute = function(el, name){
	if(typeof el.dataset !== "undefined"){
		return el.dataset[name];
	} else {
		return el.getAttribute("data-" + name);
	}
};

this.Hash = {
	save: function(){
		DocumentCookies.setItem("last-hash", document.location.hash, Infinity, "/");
	},
	restore: function(){
		var lastHash = DocumentCookies.getItem("hash");
		if(lastHash && (document.location.hash == "")){
			document.location.hash = lastHash;
		}
		DocumentCookies.removeItem("hash");
	}
};

this.GenerateID = function(){
	return Math.random().toString(16).substr(2) +
	       Math.random().toString(16).substr(2);
};