import "util/notifier.js";
import "kb.js";
import "convert.js";
import "stage.js";

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

		stageById: function(id){
			for(var i = 0; i < this.stages.length; i += 1){
				if(this.stages[i].id == id){
					return this.stages[i];
				}
			}
			return undefined;
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
			this.removeListeners();
			this.stages = [];
			this.changed();
		},

		closeLast: function(){
			// always keep one stage open
			if(this.stages.length > 1){
				this.stages[this.stages.length-1].close();
			}
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
		// title
		// link
		// after, optional
		// insteadOf, optional
		open: function(props){
			this.trim_(props.after);
			var stage = new KB.Stage(props);

			var i = this.indexOf_(props.insteadOf);
			if(i >= 0){
				this.stages[i].remove(this);
				this.stages[i] = stage;
			} else {
				this.stages.push(stage);
			}

			stage.on("closed", this.handleClose, this);
			this.changed();
			return stage.id;
		},

		openLink: function(link){
			this.open(Convert.LinkToReference(link));
		},

		handleClose: function(ev){
			this.stages = this.stages.filter(function(stage){
				return stage != ev.stage;
			});
			this.changed();
		},

		removeListeners: function(){
			this.stages.map(function(stage){
				stage.remove(this);
			});
		},
		addListeners: function(){
			var self = this;
			this.stages.map(function(stage){
				stage.on("closed", self.handleClose, self);
			});
		},

		updateRefs: function(nextstages){
			this.removeListeners();

			var stages = this.stages.slice();
			var changed = false;

			var self = this;
			var newstages = nextstages.map(function(stage){
				var prev = stages.shift();

				if(prev){
					var plink = Convert.ReferenceToLink(prev);
					var slink = Convert.ReferenceToLink(stage);
					if(plink == slink) {
						return prev;
					}
				}
				changed = true;
				return stage;
			});

			if(stages.length > 0){
				changed = true;
			}
			if(changed){
				this.stages = newstages;
				this.changed();
			}

			this.addListeners();
		},

		findStageFromElement: function(el){
			for(var i = 0; i < 32; i += 1){
				if(el == null){ return null; }
				if(el.classList.contains("stage")){
					return this.stageById(el.dataset.id);
				}
				el = el.parentElement;
			}
			return undefined;
		},

		handleOpenLink: function(ev){
			var target = ev.target;
			var stage = this.findStageFromElement(target);

			var ref = Convert.LinkToReference(target.href);
			var url = ref.url;

			if(stage){
				var locFrom = Convert.URLToLocation(stage.url);
				var locTo = Convert.URLToLocation(url);
				url = "//" + locFrom.host + locTo.pathname;
			}

			var link = target.dataset.link || ref.link;
			var title = target.innerText;

			if(ev.button == 1){
				this.open({
					url: url,
					link: link,
					title: title
				});
			} else {
				this.open({
					url: url,
					link: link,
					title: title,
					after: stage && stage.id
				});
			}

			ev.preventDefault();
		},

		handleClickLink: function(ev){
			if(ev.target.localName != "a") return;
			if(ev.target.classList.contains("external-link")) return;
			if(ev.target.onclick != null) return;
			if(ev.target.onmousedown != null) return;
			if(ev.target.onmouseup != null) return;
			if(ev.target.href == "") return;

			this.handleOpenLink(ev);
		}
	};


	function elementIsEditable(elem){
		return elem && (
			((elem.nodeName === 'INPUT') && (elem.type === 'text')) ||
			(elem.nodeName === 'TEXTAREA') ||
			(elem.contentEditable === 'true')
		);
	}

	return Lineup;
})();
