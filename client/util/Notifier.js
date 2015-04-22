'use strict';

export class Notifier {
	constructor(){
		this.listeners = [];
		this.lastKey = 0;
	}

	listen(fn, recv){
		if(typeof fn === 'undefined'){
			throw new Error("fn must be defined");
		}
		var listener = {
			fn: fn,
			recv: recv,
			key: this.lastKey++
		};
		this.listeners.push(listener);
		var self = this;
		return function(){
			self._remove(listener.key);
		};
	}
	unlisten(fn, recv){
		this.listeners = this.listeners.filter(function(listener){
			return !((listener.fn === fn) && (listener.recv === recv));
		});
	}

	_remove(key){
		this.listeners = this.listeners.filter(function(listener){
			return listener.key !== key;
		});
	}

	notify(){
		var args = arguments;
		var self = this;
		window.setTimeout(function(){
			self.update.apply(self, args);
		}, 0);
	}

	update(){
		var args = arguments;
		this.listeners.map(function(listener){
			listener.fn.apply(listener.recv, args);
		});
	}
}
