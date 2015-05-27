function FindListContainer(node, containerClass, listClass){
	while(node){
		if(node.classList.contains(listClass)){
			return node;
		}
		if(node.classList.contains(containerClass)){
			return node.getElementsByClassName(listClass)[0];
		}
		node = node.parentElement;
	}
	return null;
}

function FindListPosition(ev, containerClass, listClass, skip){
	var list = FindListContainer(ev.target, containerClass, listClass);
	if(list == null){
		return;
	}

	var mouse = {
		x: ev.pageX,
		y: ev.pageY
	};

	var child = list.firstChild;
	if(child){
		var box = child.getBoundingClientRect();
		if(child && (mouse.y < box.top)){
			return {node: null, rel: "before"};
		}
	} else {
		return {node: null, rel: "before"};
	}

	for(var i = 0; i < list.children.length; i += 1){
		var child = list.children[i];
		if(skip && skip(child)){
			continue;
		}

		var box = child.getBoundingClientRect();
		if(mouse.y > box.bottom){
			continue;
		}
		return {node: child, rel: "after"};
	}

	return {node: null, rel: "after"};
}
