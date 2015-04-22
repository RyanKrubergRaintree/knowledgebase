'use strict';

import {convert} from './convert';

export class Page {
	constructor(page){
		var page = page || {};

		this.meta = page.meta || { upvotes: 0, downvotes: 0, tags: [] };
		if(page.meta){
			this.meta.upvotes = page.meta.upvotes || 0;
			this.meta.downvotes = page.meta.downvotes || 0;
			this.meta.tags = this.meta.tags || [];
		}


		this.flags = page.flags ||[];
		this.path = page.path || '';
		this.url = page.url || this.path;
		this.title = page.title || '';
		this.synopsis = page.synopsis || '';
		this.version = page.version || 0;

		this.comments = page.comments || [];
		this.story = page.story || [];
		this.journal = page.journal || [];

		var now = new Date();
		this.created = page.created || now.valueOf();
		this.modified = page.modified || now.valueOf();
	}

	//TODO: remove these
	get readonly(){
		return this.flags.indexOf('read-only') >= 0;
	}
	get dynamic(){
		return this.flags.indexOf('dynamic') >= 0;
	}
	get editableHeader(){
		return this.flags.indexOf('editable-header') >= 0;
	}

	updateURL(responseURL){
		this.url = responseURL;
	}

	indexOf(itemId){
		var story = this.story;
		for(var i = 0; i < story.length; i += 1){
			if(story[i].id == itemId){
				return i;
			}
		}
		throw 'Item "' + itemId + '" does not exist.';
	}
	apply(op){
		var fn = ops[op.type];
		if(fn){
			fn(this, this.story, op);
			this.version += 1;
			this.journal.push(op);
			this.modified = (new Date()).valueOf();
			return;
		}
		throw new Error('unknown op type ' + op.type);
	}
}

var ops = {}

// page operations
ops['add'] = function(page, story, op){
	if(op.after != null){
		var i = page.indexOf(op.after);
		story.splice(i+1, 0, op.item);
	} else {
		story.unshift(op.item);
	}
};

ops['remove'] = function(page, story, op){
	var i = page.indexOf(op.id);
	story.splice(i, 1);
};

ops['edit'] = function(page, story, op){
	var i = page.indexOf(op.id);
	story[i] = op.item;
};

ops['edit-text'] = function(page, story, op){
	var i = page.indexOf(op.id);
	story[i].text = op.item.text;
};

ops['move'] = function(page, story, op){
	var from = page.indexOf(op.id),
		item = story.splice(from, 1)[0];
	if(op.after != null){
		var to = page.indexOf(op.after);
		story.splice(to+1, 0, item);
	} else {
		story.unshift(item);
	}
};

ops['header'] = function(page, story, op){
	page.tags = op.tags;
	page.synopsis = op.synopsis;
};

ops['vote-up'] = function(page, story, op){
	page.meta.upvotes = (page.meta.upvotes || 0) + 1;
};

ops['vote-down'] = function(page, story, op){
	page.meta.downvotes = (page.meta.downvotes || 0) + 1;
};

ops['comment-add'] = function(page, story, op){
	page.comments.push(op.comment);
};

ops['comment-remove'] = function(page, story, op){
	var cs = page.comments;
	for(var i = 0; i < cs.length; i += 1){
		if(cs[i].id == op.id){
			cs.splice(i, 1);
			return;
		}
	}
};
