package('kb', function(exports){
	'use strict';

	depends('Lineup.View.js');

	var HeaderMenu = React.createClass({
		displayName: 'HeaderMenu',
		render: function(){
			return React.DOM.div({
				className:'header-menu'
			},
				this.props.items.map(function(item){
					return React.DOM.a(item, item.caption);
				})
			);
		}
	});

	var Search = React.createClass({
		lastStageId: undefined,
		search: function(ev){
			var Lineup = this.props.Lineup;
			var query = this.refs.query.getDOMNode().value.trim();
			if(ev.shiftKey){
				this.lastStageId = undefined;
			}
			this.lastStageId = Lineup.open({
				url: '/search:search?q='+query,
				title: '"' + query + '"',
				insteadOf: this.lastStageId
			});
			ev.preventDefault();
		},
		keyDown: function(ev){
			var Lineup = this.props.Lineup;
			if(ev.keyCode === 27){// esc
				Lineup.closeLast();
				return;
			}

			if(ev.keyCode === 13){
				// open page directly
				if(ev.ctrlKey){
					if(!ev.shiftKey){
						Lineup.clear();
					}

					var query = this.refs.query.getDOMNode().value.trim();
					Lineup.openLink(query);
					ev.preventDefault();
					return;
				}

				this.search(ev);
				ev.preventDefault();
				return;
			}

			var stages = document.querySelectorAll('.stage');
			if(stages.length === 0){
				return;
			}

			var stage = stages[stages.length-1];
			var middle = stage.querySelector('.stage-scroll');

			switch(ev.keyCode){
			case 33: // pageup
				middle.scrollTop -= middle.clientHeight;
				break;
			case 34: // pagedown
				middle.scrollTop += middle.clientHeight;
				break;
			}
		},
		render: function(){
			return React.DOM.form(
				{
					className:'search',
					onSubmit: this.search
				},
				React.DOM.input({
					ref: 'query',
					placeholder:'Search...',
					onKeyDown: this.keyDown
				}),
				React.DOM.button({
					className:'search-icon mdi mdi-magnify',
					type: 'submit',
					tabIndex: -1
				})
			);
		}
	});

	var Header = React.createClass({
		openHome: function(ev){
			ev.preventDefault();
			ev.stopPropagation();

			var lineup = this.props.Lineup;
			lineup.clear();
			lineup.openLink(KBHomePage);
		},
		createNewPage: function(ev){
			ev.preventDefault();
			ev.stopPropagation();

			var lineup = this.props.Lineup;
			lineup.open({
				url: '',
				link: '',
				title: ''
			});
		},
		logout: function(ev){
			ev.preventDefault();
			ev.stopPropagation();

			window.location.pathname = '/system/auth/logout';
		},
		displayName: 'Header',
		render: function(){
			var a = React.DOM.a;
			return React.DOM.div({
				id:'header'
			},
				a({className:'button home mdi mdi-home', href:'#', title:'Home', onClick: this.openHome}),
				React.createElement(Search, this.props),
				a({
					className:'button userinfo',
					id:'userinfo',
					href:'/user:current'
				},
					KBUser
				),
				React.createElement(HeaderMenu, {
					items: [
						{key:'0', href: '#', onClick: this.createNewPage, caption: 'New Page'},
						{key:'1', href: '/page:recent-changes', caption: 'Recent Changes'},
						{key:'2', href: '#', onClick: this.logout, caption: 'Logout'}
					]
				})
			);
		}
	});

	var Content = React.createClass({
		displayName: 'Content',
		render: function(){
			return React.DOM.div({
				id: 'content'
			},
				React.createElement(kb.Lineup.View, this.props)
			);
		}
	});

	exports.Site = React.createClass({
		displayName: 'Site',

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
});