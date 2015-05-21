//import "/kb/KB.js"

KB.Page = (function(){
	"use strict";

	// Page corresponds to the knowledgebase page structure
	function Page(data){
		data = data || {};
		this.owner = data.owner || "";
		this.slug = data.slug || "";
		this.title = data.title || "";
		this.synopsis = data.synopsis || "";
		this.story = data.story || [
			{id:"1", type:"paragraph", text: "Quibusdam distinctio, sapiente beatae laudantium ipsa quisquam. Minima nostrum non tenetur cumque, ipsum perspiciatis animi illum! Aperiam nihil, culpa laboriosam voluptatibus ratione."},
			{id:"2", type:"paragraph", text: "Consequuntur reprehenderit sint id, assumenda quisquam! Quo, libero. Assumenda, sint aspernatur! Tenetur quos cupiditate adipisci omnis laboriosam sequi labore praesentium molestias iste!"},
			{id:"3", type:"paragraph", text: "Placeat nemo sed dolore, neque, voluptatum quam, eveniet delectus, in non itaque corporis praesentium similique perspiciatis illum odit culpa assumenda adipisci. Voluptatibus."},
			{id:"4", type:"paragraph", text: "Ipsam doloribus culpa, autem quas illum velit ut aut adipisci error. Omnis sint laborum quisquam hic debitis numquam atque eaque fugit ipsam?"},
			{id:"5", type:"paragraph", text: "Nihil totam, alias accusantium quisquam ad nisi fuga deserunt labore placeat nesciunt deleniti earum eos magnam, debitis, minus non sapiente exercitationem eveniet."}
		];
		this.journal = data.journal || [];
		this.version = data.version || 0;
		//TODO: use correct date format
		this.modified = data.modified || (new Date()).valueOf();
	};

	Page.prototype = {
		clone: function(){
			return new Page(JSON.stringify(this));
		},

		indexOf_: function(id){
			var story = this.story;
			for(var i = 0; i < story.length; i += 1){
				if(story[i].id == itemId){
					return i;
				}
			}
			throw new Error('Item "' + itemId + '" does not exist.');
		},

		apply: function(op){
			var fn = OP[op.type];
			if(fn){
				fn(this, this.story, op);
				this.version += 1;
				this.journal.push(op);
				this.modified = (new Date()).valueOf();
				return;
			}
			throw new Error('Unknown operation "' + op.type + '"');
		}
	};

	var OP = {};
	OP.add = function(page, stroy, op){
		if(op.after != null){
			var i = page.indexOf_(op.after);
			story.splice(i+1, 0, op.item);
		} else {
			story.unshift(op.item);
		}
	};

	OP.remove = function(page, stroy, op){
		var i = page.indexOf_(op.id);
		story.splice(i, 1);
	};

	OP.edit = function(page, stroy, op){
		var i = page.indexOf_(op.id);
		story[i] = op.item;
	};

	OP.move = function(page, stroy, op){
		var from = page.indexOf_(op.id),
			item = story.splice(from, 1)[0];
		if(op.after != null){
			var to = page.indexOf_(op.after);
			story.splice(to+1, 0, item);
		} else {
			story.unshift(item);
		}
	};

	return Page;
})();
