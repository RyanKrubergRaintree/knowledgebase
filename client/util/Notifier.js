"use strict";

function Notifier(){
	this.listeners = [];
}

Notifier.prototype = {
	mixto: function(obj){
		var self = this;
		obj.on = function(event, handler, recv){
			self.on(event, handler, recv);
		};
		obj.off = function(event, handler, recv){
			self.off(event, handler, recv);
		};
		obj.remove = function(recv){
			self.remove(recv);
		};
	},
	on: function(event, handler, recv){
		this.listeners.push({
			event: event,
			handler: handler,
			recv: recv
		});
	},
	off: function(event, handler, recv){
		this.listeners = this.listeners.filter(
			function(listener){
				return !(
					(listener.event === event) &&
					(listener.handler === handler) &&
					(listener.recv === recv)
				);
			}
		);
	},
	remove: function(recv){
		this.listeners = this.listeners.filter(
			function(listener){
				return !(listener.recv === recv);
			}
		);
	},
	emit: function(event){
		var self = this;
		window.setTimeout(function(){
			self.handle(event);
		}, 0);
	},
	handle: function(event){
		this.listeners.map(function(listener){
			if(listener.event == event.type){
				listener.handler.call(listener.recv, event);
			}
		});
	}
};
