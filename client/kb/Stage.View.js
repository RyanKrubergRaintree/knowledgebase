package('kb.Stage', function(exports) {
	'use strict';

	depends('util/SmoothScroll.js');

	depends('Convert.js');
	depends('Stage.js');
	depends('Page.View.js');

	var StageButtons = React.createClass({
		displayName: 'StageButtons',

		toggleWidth: function() {
			this.props.onToggleWidth();
		},
		close: function() {
			this.props.stage.close();
		},

		createFactory: function(ev) {
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

		deletePage: function() {
			var stage = this.props.stage;
			var check = window.prompt('Delete this page?\nType the page link \'' + stage.link + '\' to confirm:');
			if (check === null) {
				return;
			}
			if (check.trim() !== stage.link.trim()) {
				return;
			}

			stage.destroy();
		},

		render: function() {
			var stage = this.props.stage;
			var a = React.DOM.a;
			return React.DOM.div({
					className: 'stage-buttons'
				},
				stage.canModify() ? React.DOM.a({
					className: 'mdi mdi-playlist-plus',
					title: 'Drag to page to add an item.',
					style: {
						cursor: 'move'
					},
					draggable: 'true',
					href: '#',
					onClick: function(ev) {
						ev = ev || window.event;
						ev.preventDefault();
					},
					onDragStart: this.createFactory
				}) : null,
				stage.canDestroy() ? a({
					className: 'mdi mdi-delete',
					title: 'Delete this page.',
					onClick: this.deletePage
				}) : null,
				a({
					className: 'mdi ' + (this.props.isWide ? 'mdi-arrow-collapse' : 'mdi-arrow-expand'),
					title: 'Toggle page width.',
					onClick: this.toggleWidth
				}),
				a({
					className: 'mdi mdi-close',
					title: 'Close page.',
					onClick: this.close
				})
			);
		}
	});

	var NewPage = React.createClass({
		displayName: 'NewPage',
		tryCreate: function(ev) {
			var stage = this.props.stage;

			stage.title = this.state.title;
			stage.link = kb.convert.TextToSlug(this.state.owner + ':' + stage.title);
			stage.create();

			ev.preventDefault();
			ev.stopPropagation();
		},

		getInitialState: function() {
			return {
				title: this.props.stage.title,
				owner: kb.convert.LinkToOwner(this.props.stage.link) || '',
				groups: []
			};
		},

		groupsReceived: function(xhr) {
			if (xhr.status === 200) {
				var info = JSON.parse(xhr.responseText);
				this.setState({
					groups: info.groups || []
				});
			}
		},

		componentDidMount: function() {
			var xhr = new XMLHttpRequest();
			var self = this;
			xhr.onreadystatechange = function() {
				if (xhr.readyState !== 4) {
					return;
				}
				self.groupsReceived(xhr);
			};
			xhr.open('GET', '/user:editor-groups', true);
			xhr.setRequestHeader('Accept', 'application/json');
			xhr.send();
		},

		ownerChanged: function(ev) {
			this.setState({
				owner: ev.currentTarget.value
			});
		},

		titleChanged: function() {
			this.setState({
				title: this.refs.title.getDOMNode().value
			});
		},

		render: function() {
			var self = this;
			var stage = this.props.stage;
			var title = this.state.title,
				owner = this.state.owner,
				link = kb.convert.TextToSlug(owner + ':' + title);

			return React.DOM.div({
					className: 'page new-page'
				},
				React.DOM.form({
						onSubmit: this.tryCreate
					},
					React.DOM.label({}, 'Link'),
					React.DOM.span({
						className: 'link'
					}, link),
					React.DOM.label({
						htmlFor: 'new-page-title'
					}, 'Title'),
					React.DOM.input({
						id: 'new-page-title',
						className: 'title',
						ref: 'title',
						defaultValue: stage.title,
						onChange: this.titleChanged,
						onKeyUp: this.titleChanged,
						autoFocus: true
					}),
					React.DOM.label({}, 'Owner'),
					React.DOM.div({
							className: 'group'
						},
						this.state.groups.map(function(group) {
							var checked = owner === group;
							return (
								React.DOM.div({
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
					React.DOM.button({
						type: 'submit'
					}, 'Create')
				)
			);
		}
	});

	exports.View = React.createClass({
		displayName: 'Stage',

		toggleWidth: function() {
			if (this.props.stage.wide) {
				this.props.stage.collapse();
			} else {
				this.props.stage.expand();
			}

			window.setTimeout(this.activate, 100);
		},
		activate: function(ev) {
			if (typeof ev === 'undefined') {
				kb.util.SmoothScroll.to(this.getDOMNode());
			} else if (!ev.defaultPrevented) {
				kb.util.SmoothScroll.to(this.getDOMNode());
			}
		},

		render: function() {
			var stage = this.props.stage;
			if (stage.creating) {
				return React.DOM.div({
						className: 'stage',
						onClick: this.activate,
						'data-id': stage.id,

						style: this.props.style
					},
					React.createElement(StageButtons, {
						stage: this.props.stage,
						isWide: stage.wide,
						onToggleWidth: this.toggleWidth
					}),
					React.DOM.div({
							className: 'stage-scroll round-scrollbar'
						},
						React.createElement(NewPage, {
							stage: this.props.stage
						})
					)
				);
			}

			return React.DOM.div({
					className: 'stage',
					onClick: this.activate,
					'data-id': stage.id,

					style: this.props.style
				},
				React.createElement(StageButtons, {
					stage: this.props.stage,
					isWide: stage.wide,
					onToggleWidth: this.toggleWidth
				}),
				React.DOM.div({
						className: 'stage-scroll round-scrollbar'
					},
					React.createElement(kb.Page.View, {
						stage: this.props.stage,
						page: this.props.stage.page
					})
				)
			);
		},

		// bindings to Stage
		changed: function() {
			this.forceUpdate();
		},
		widthChanged: function() {
			this.props.onWidthChanged();
		},
		componentDidMount: function() {
			this.props.stage.on('changed', this.changed, this);
			this.props.stage.on('widthChanged', this.widthChanged, this);
			this.props.stage.pull();
			this.activate();
		},
		componentWillReceiveProps: function(nextprops) {
			if (this.props.stage !== nextprops.stage) {
				this.props.stage.remove(this);
				nextprops.stage.on('changed', this.changed, this);
				nextprops.stage.on('widthChanged', this.widthChanged, this);
				nextprops.stage.pull();
			}
		},
		componentWillUnmount: function() {
			this.props.stage.remove(this);
		}
	});
});
