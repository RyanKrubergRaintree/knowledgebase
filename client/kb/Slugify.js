package('kb', function(exports){
	'use strict';

	depends('unicode/identifier.js');
	depends('unicode/runename.js');

	// Slugify converts text to a slug
	//
	// * numbers and '/' are left intact
	// * letters will be lowercased (if possible)
	// * '-', ',', '.', ' ', '_' will be converted to '-'
	// * other symbols or punctuations will be converted to html entity reference name
	//   (if there exists such reference name)
	// * everything else will be converted to '-'
	//
	// Example:
	//   '&Hello_世界/+!' ==> 'amp-hello-世界/plus-excl'
	//   'Hello  World  /  Test' ==> 'hello-world/test'
	exports.Slugify = Slugify;
	function Slugify(title){
		var cutdash = true,
			emitdash = false;

		var slug = '';

		for(var i = 0; i < title.length; i += 1){
			var r = title[i];
			if(kb.unicode.IsIdent(r)) {
				if(emitdash && !cutdash){
					slug += '-';
				}
				slug += r.toLowerCase();

				emitdash = false;
				cutdash = false;
				continue;
			}
			if((r == '/') || (r == ':')){
				slug += r;
				emitdash = false;
				cutdash = true;
			} else if ((r == '-') || (r == ',') || (r == '.') || (r == ' ') || (r == '_')) {
				emitdash = true;
			} else {
				var name = kb.unicode.RuneName[r];
				if(name){
					if(!cutdash){
						slug += '-';
					}
					slug += name;
					cutdash = false;
				}
				emitdash = true;
			}
		}

		if(slug.length == 0){
			return '-';
		}

		return slug;
	}
});