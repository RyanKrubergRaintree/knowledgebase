function Resolve(stage, text){
	if(text == null) {
		return '';
	}

	text = text.replace(Resolve.rxExternalLink,
		'<a href="$1" class="external-link" target="_blank" rel="nofollow">$2</a>');

	text = text.replace(Resolve.rxInternalLink, function(match, link){
		var ref = Convert.LinkToReference(link, stage);
		return '<a href="' + ref.url + '" data-link="' + ref.link + '" >' + ref.title + "</a>";
	});
	return text;
}
Resolve.rxExternalLink = /\[\[\s*(https?\:[^ \]]+)\s+([^\]]+)\]\]/g;
Resolve.rxInternalLink = /\[\[\s*([^\]]+)\s*\]\]/g;

//TODO: add HTML sanitation
function ResolveHTML(stage, text){
	return Resolve(stage, text);
}
