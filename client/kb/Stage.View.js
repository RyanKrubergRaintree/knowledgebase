package('kb.Stage', function(exports){
	'use strict';

	depends('util/SmoothScroll.js');

	depends('Convert.js');
	depends('Stage.js');
	depends('Page.View.js');

	var StageButtons = React.createClass({
		displayName: 'StageButtons',

		toggleWidth: function(){
			this.props.onToggleWidth();
		},
		close: function(){
			this.props.stage.close();
		},

		createFactory: function(ev){
			var item = {
				id: GenerateID(),
				type: 'factory',
				text: ''
			};

			ev.dataTransfer.effectAllowed = 'copy';
			var data = {
				item: item
			};

			ev.dataTransfer.setData('Text', JSON.stringify(data));
		},

		deletePage: function(){
			var stage = this.props.stage;
			var check = window.prompt('Delete this page?\nType the page link \''+ stage.link + '\' to confirm:');
			if(check.trim() !== stage.link.trim()){
				return;
			}

			stage.destroy();
		},

		render: function(){
			var stage = this.props.stage;
			var a = React.DOM.a;
			return React.DOM.div(
				{className: 'stage-buttons'},
				stage.canModify() ? React.DOM.div({
					className:'mdi mdi-playlist-plus',
					title:'Drag to page to add an item.',
					style: { cursor: 'move' },
					draggable: 'true',
					onDragStart: this.createFactory
				}) : null,
				stage.canDestroy() ? a({
					className:'mdi mdi-delete',
					title:'Delete this page.',
					onClick: this.deletePage
				}) : null,
				a({
					className:'mdi ' + (this.props.isWide ? 'mdi-arrow-collapse' : 'mdi-arrow-expand'),
					title:'Toggle page width.',
					onClick: this.toggleWidth
				}),
				a({
					className:'mdi mdi-close',
					title:'Close page.',
					onClick: this.close
				})
			);
		}
	});

	var NewPage = React.createClass({
		displayName: 'NewPage',
		tryCreate: function(ev){
			var stage = this.props.stage;

			stage.title = this.state.title;
			stage.link  = kb.convert.TextToSlug(this.state.owner + ':' + stage.title);
			stage.create();

			ev.preventDefault();
			ev.stopPropagation();
		},

		getInitialState: function(){
			return {
				title: this.props.stage.title,
				owner: kb.convert.LinkToOwner(this.props.stage.link) || '',
				groups: []
			};
		},

		groupsReceived: function(ev){
			var xhr = ev.target;
			if(xhr.status === 200){
				var info = JSON.parse(ev.target.responseText);
				this.setState({groups: info.groups || []});
			}
		},

		componentDidMount: function(){
			var xhr = new XMLHttpRequest();
			xhr.onload = this.groupsReceived;
			xhr.open('RAW', '/user:editor-groups', true);
			xhr.setRequestHeader('Accept', 'application/json');
			xhr.send();
		},

		ownerChanged: function(ev){
			this.setState({
				owner: ev.currentTarget.value
			});
		},

		titleChanged: function(){
			this.setState({
				title: this.refs.title.getDOMNode().value
			});
		},

		render: function(){
			var self = this;
			var stage = this.props.stage;
			var title = this.state.title,
				owner = this.state.owner,
				link = kb.convert.TextToSlug(owner + ':' + title);

			return React.DOM.div(
				{ className: 'page new-page' },
				React.DOM.form({
					onSubmit: this.tryCreate
				},
					React.DOM.label({}, 'Link'),
					React.DOM.span({className:'link'}, link),
					React.DOM.label({ htmlFor: 'new-page-title' }, 'Title'),
					React.DOM.input({
						id: 'new-page-title',
						className: 'title',
						ref: 'title',
						defaultValue: stage.title,
						onChange: this.titleChanged,
						autoFocus: true
					}),
					React.DOM.label({}, 'Owner'),
					React.DOM.div(
						{ className: 'group' },
						this.state.groups.map(function(group){
							var checked = owner === group;
							return (
								React.DOM.div(
									{
										key: group,
										className: checked ? 'checked' : ''
									},
									React.DOM.input({
										id: 'group-' + group,
										type: 'radio',
										name: 'group',
										value: group,
										onChange: self.ownerChanged,
										checked: checked
									}),
									React.DOM.label({
										htmlFor: 'group-' + group
									}, group)
								)
							);
						})
					),
					React.DOM.button({ type: 'submit' }, 'Create')
				)
			);
		}
	});

	exports.View = React.createClass({
		displayName: 'Stage',

		getInitialState: function(){
			return {
				wide: false
			};
		},

		toggleWidth: function(){
			this.setState({
				wide: !this.state.wide
			});
			window.setTimeout(this.activate, 100);
		},
		activate: function(ev){
			if(typeof ev === 'undefined'){
				kb.util.SmoothScroll.to(this.getDOMNode());
			} else if (!ev.defaultPrevented){
				kb.util.SmoothScroll.to(this.getDOMNode());
			}
		},

		render: function(){
			var wide = this.state.wide ? ' stage-wide' : '';
			var stage = this.props.stage;

			if(stage.creating){
				return React.DOM.div(
					{
						className: 'stage' + wide,
						onClick: this.activate,
						'data-id': stage.id
					},
					React.createElement(StageButtons, {
						stage: this.props.stage,
						isWide: this.state.wide,
						onToggleWidth: this.toggleWidth
					}),
					React.DOM.div(
						{ className:'stage-scroll round-scrollbar'},
						React.createElement(NewPage, {
							stage: this.props.stage
						})
					)
				);
			}

			return React.DOM.div(
				{
					className: 'stage' + wide,
					onClick: this.activate,
					'data-id': stage.id
				},
				React.createElement(StageButtons, {
					stage: this.props.stage,
					isWide: this.state.wide,
					onToggleWidth: this.toggleWidth
				}),
				React.DOM.div(
					{ className:'stage-scroll round-scrollbar'},
					React.createElement(kb.Page.View, {
						stage: this.props.stage,
						page: this.props.stage.page
					})
				)
			);
		},

		// bindings to Stage
		changed: function(){
			this.forceUpdate();
		},
		componentDidMount: function(){
			this.props.stage.on('changed', this.changed, this);
			this.props.stage.pull();
			this.activate();
		},
		componentWillReceiveProps: function(nextprops){
			if(this.props.stage !== nextprops.stage){
				this.props.stage.remove(this);
				nextprops.stage.on('changed', this.changed, this);
				nextprops.stage.pull();
			}
		},
		componentWillUnmount: function() {
			this.props.stage.remove(this);
		}
	});
});