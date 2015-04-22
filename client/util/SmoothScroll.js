var timer = null,
	speed = 0.09;

// returns true, if no additional movement is necessary
function scrollTowards(target){
	if((target == null) || (target.parentNode == null)){
		return true;
	}
	var parent = target.parentNode,
		overTop = target.offsetTop - parent.offsetTop < parent.scrollTop,
		overBottom = (target.offsetTop - parent.offsetTop + target.clientHeight) > (parent.scrollTop + parent.clientHeight),
		overLeft = target.offsetLeft - parent.offsetLeft < parent.scrollLeft,
		overRight = (target.offsetLeft - parent.offsetLeft + target.clientWidth) > (parent.scrollLeft + parent.clientWidth),
		alignWithTop = overTop && !overBottom;

	var prevTop = parent.scrollTop;
	var prevLeft = parent.scrollLeft;

	var newTop = parent.scrollTop;
	if (overTop || overBottom) {
		newTop = target.offsetTop - parent.offsetTop - parent.clientHeight / 2 + target.clientHeight / 2;
	}
	parent.scrollTop = (speed*newTop + (1 - speed)*parent.scrollTop);
	if(Math.abs(parent.scrollTop - newTop) < 2){
		parent.scrollTop = newTop;
	}

	var newLeft = parent.scrollLeft;
	if (overLeft || overRight) {
		newLeft = target.offsetLeft - parent.offsetLeft - parent.clientWidth / 2 + target.clientWidth / 2;
	}
	parent.scrollLeft = (speed*newLeft + (1 - speed)*parent.scrollLeft);
	if(Math.abs(parent.scrollLeft - newLeft) < 2){
		parent.scrollLeft = newLeft;
	}

	if((parent.scrollTop == prevTop) && (parent.scrollLeft == prevLeft)){
		return true;
	}
	return (parent.scrollTop == newTop) && (parent.scrollLeft == newLeft)
};

function cancel(){
	if(timer){
		if(window.requestAnimationFrame){
			window.cancelAnimationFrame(timer);
		} else {
			window.clearTimeout(timer);
		}
		timer = null;
	}
}

function to(target){
	cancel();
	if(scrollTowards(target)){
		return;
	};

	if(window.requestAnimationFrame){
		timer = window.requestAnimationFrame(function(){ to(target); });
	} else {
		timer = window.setTimeout(function(){ to(target); }, 15);
	}
}

export var SmoothScroll = {
	to: to
};
