import "kb.js"
import "lineup.view.js"

KB.Site = (function(){
	"use strict"

	var HeaderMenu = React.createClass({
		displayName: "HeaderMenu",
		render: function(){
			return React.DOM.div({
				className:"header-menu"
			},
				this.props.items.map(function(item, i){
					return React.DOM.a(item, item.caption);
				})
			);
		}
	})

	var Search = React.createClass({
		lastStageId: undefined,
		search: function(ev){
			var Lineup = this.props.Lineup;
			var query = this.refs.query.getDOMNode().value.trim();
			this.lastStageId = Lineup.open({
				url: "/page:search?q="+query,
				title: 'Search "' + query + '"',
				insteadOf: this.lastStageId
			});
			ev.preventDefault();
		},
		keyDown: function(ev){
			var Lineup = this.props.Lineup;
			if(ev.keyCode == 27){// esc
				Lineup.closeLast();
				return;
			}

			if((ev.keyCode == 13) && (ev.ctrlKey || ev.shiftKey)){
				if(ev.ctrlKey){
					Lineup.clear();
				}

				var query = this.refs.query.getDOMNode().value.trim();
				Lineup.openLink(query);
				ev.preventDefault();
				return;
			}

			var stages = document.getElementsByClassName('stage');
			if(stages.length == 0){
				return;
			};

			var stage = stages[stages.length-1];
			var middle = stage.getElementsByClassName('stage-scroll')[0];

			switch(ev.keyCode){
			case 33: // pageup
				middle.scrollTop -= middle.clientHeight
				break;
			case 34: // pagedown
				middle.scrollTop += middle.clientHeight
				break;
			}
		},
		render: function(){
			return React.DOM.form(
				{
					className:"search",
					onSubmit: this.search
				},
				React.DOM.input({
					ref: "query",
					placeholder:"Search...",
					onKeyDown: this.keyDown
				}),
				React.DOM.button({
					className:"search-icon mdi mdi-magnify",
					type: "submit",
					tabIndex: -1
				})
			);
		}
	})

	var Header = React.createClass({
		openHome: function(ev){
			var lineup = this.props.Lineup;
			lineup.clear();
			lineup.openLink(Global.HomePage);
			ev.preventDefault();
			ev.stopPropagation();
		},
		createNewPage: function(ev){
			var lineup = this.props.Lineup;
			lineup.open({
				url: null,
				link: "",
				title: ""
			});
			ev.preventDefault();
			ev.stopPropagation();
		},
		displayName: "Header",
		render: function(){
			var a = React.DOM.a;
			return React.DOM.div({
				id:"header"
			},
				a({className:"button home mdi mdi-home", href:"#", title:"Home", onClick: this.openHome}),
				React.createElement(Search, this.props),
				a({
					className:"button userinfo",
					id:"userinfo",
					href:"/user:current"
				},
					Global.User
				),
				React.createElement(HeaderMenu, {
					items: [
						{key:"0", href: "#", onClick: this.createNewPage, caption: "New Page"},
						{key:"1", href: "/page:recent-changes", caption: "Recent Changes"},
						{key:"2", href: "/system/auth/logout", caption: "Logout"}
					]
				})
			);
		}
	});

	var Content = React.createClass({
		displayName: "Content",
		render: function(){
			return React.DOM.div({
				id: "content"
			},
				React.createElement(KB.Lineup.View, this.props)
			);
		}
	});

	var Site = React.createClass({
		displayName: "Site",

		componentDidMount: function(){
			var self = this;
			window.LiveBundleChange = function(){
				self.forceUpdate();
			};
		},

		componentWillUnmount: function(){
			window.LiveBundleChange = function(){};
		},

		render: function(){
			return React.DOM.div({},
				React.createElement(Header,  this.props),
				React.createElement(Content, this.props)
			);
		}
	});

	return Site;
})();
