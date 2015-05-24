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

	var Header = React.createClass({
		openHome: function(ev){
			var lineup = this.props.Lineup;
			lineup.clear();
			lineup.openLink(Global.HomePage);
			ev.preventDefault();
			ev.stopPropagation();
		},
		displayName: "Header",
		render: function(){
			var a = React.DOM.a;
			return React.DOM.div({
				id:"header"
			},
				a({className:"button logo", href:"#", title:"Home", onClick: this.openHome}),
				React.DOM.form({className:"search"},
					React.DOM.input({placeholder:"Search..."}),
					React.DOM.button({
						className:"search-icon mdi mdi-magnify",
						type: "submit",
						tabIndex: -1
					})),
				a({
					className:"button userinfo",
					id:"userinfo",
					href:"/user:"+ Slugify(Global.User)
				},
					Global.User
				),
				React.createElement(HeaderMenu, {
					items: [
						{key:"0", href: "#", caption: "New Page"},
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
