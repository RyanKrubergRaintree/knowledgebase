import "slug.js";

function TrackPageView(url, title){
	if(typeof ga !== 'undefined'){
		ga('set', {page: url, title: title});
		ga('send', 'pageview');
	}
}
