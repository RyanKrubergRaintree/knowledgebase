// import "/util/Notifier.js"
// import "/kb/KB.js"
// import "/kb/Convert.js"
// import "/kb/Stage.js"


//TODO: get rid of close --> use notification from Stage

KB.Lineup = (function(){
	"use strict";

	function Lineup(){
		this.stages = [];
		this.notifier = new Notifier();
		this.notifier.mixto(this);
	}

	Lineup.prototype = {
		changed: function(){
			this.notifier.emit({
				type:"changed",
				lineup: this
			});
		},

		indexOf_: function(id){
			if(typeof id === 'undefined'){
				return -1;
			}
			for(var i = 0; i < this.stages.length; i += 1){
				if(this.stages[i].id == id){
					return i;
				}
			}
			return -1;
		},

		trim_: function(id){
			if(typeof id === 'undefined'){
				return;
			}
			var i = this.indexOf_(id);
			if(i >= 0){
				this.stages = this.stages.slice(0, i + 1);
			}
		},

		clear: function(){
			this.stages = [];
			this.changed();
		},

		close: function(id){
			this.stages = this.stages.filter(function(stage){
				return stage.id !== id;
			});
			this.changed();
		},

		closeLast: function(){
			var stages = this.stages;
			stages = stages.slice(0, Math.max(stages.length-1, 1));
			this.changed();
		},

		changeRef: function(id, stage){
			var i = this.indexOf_(id);
			if(i >= 0){
				var ref = this.stages[i];
				ref.url = Convert.URLToReadable(stage.url);
				ref.link = Convert.URLToLink(stage.link);
				ref.title = stage.title;
				this.changed();
			}
		},


		// url
		// title, optional
		// link, optional
		// after, optional
		// insteadOf, optional
		open: function(props){
			this.trim_(props.after);
			var url = Convert.URLToReadable(props.url);

			var stage = new KB.Stage({
				url: url,
				title: props.title || Convert.URLToTitle(url),
				link: props.link || Convert.URLToLink(url)
			});

			if(props.link === ""){
				stage.link = "";
			}

			var i = this.indexOf_(props.insteadOf);
			if(i >= 0){
				this.stages[i] = stage;
			} else {
				this.stages.push(stage);
			}

			this.changed();
			return stage.id;
		},


		updateRefs: function(nextproxies){
			var preproxies = this.stages.slice();
			var changed = false;

			var self = this;
			var newproxies = nextproxies.map(function(stage){
				var prev = preproxies.shift();
				if(prev && (prev.url == stage.url)){
					return prev;
				}
				changed = true;
				return stage;
			});

			if(preproxies.length > 0){
				changed = true;
			}
			if(changed){
				this.stages = newproxies;
				this.changed();
			}
		}
	};

	return Lineup;
})();
