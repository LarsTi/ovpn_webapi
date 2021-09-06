'use strict';

class AccessGroupLine extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			subnet: '',
			name: '',
			mask: '',
			editable: false
		};
		this.modify = this.modify.bind(this);
		this.delete = this.delete.bind(this);
	}
	modify(e) {
		e.preventDefault();
		if (this.state.editable){
			if ( this.state.name === this.props.accessGroup.name
				&& this.state.subnet === this.props.accessGroup.subnet
				&& this.state.mask === this.props.accessGroup.mask ){
				//Es gab keine Änderung.
				this.setState({editable: false });
			}else{
				//Updaten der Änderung ins Backend
				// this will update entries with PUT
				fetch(window.location.toString() + "api/accessgroup/" + this.props.accessGroup.ID, {
					"method": "PUT",
					"headers": {
						"content-type": "application/json",
						"accept": "application/json"
					},
					"body": JSON.stringify({
						name: this.state.name,
						subnet: this.state.subnet,
						mask: this.state.mask
					})
				})
					.then(response => response.json())
					.then(response => { console.log(response);
					})
					.catch(err => { console.log(err); });
				this.setState({editable: false });
			}
		}else{
			this.setState({
				name:this.props.accessGroup.name,
				subnet: this.props.accessGroup.subnet,
				mask: this.props.accessGroup.mask,
				editable: true
			})
		}
	}
	delete(e) {
		// delete entity - DELETE
		e.preventDefault();
		// deletes entities
		fetch(window.location.toString() + "api/accessgroup/" + this.props.accessGroup.ID, {
			"method": "DELETE",
			"headers": {
			}
		})
			.then(response => response.json())
			.then(response => {
				console.log(response);
			})
			.catch(err => {
				console.log(err);
			})
			.finally(() => this.props.refresh());
	}

	render() {
		var _that = this;
		return React.createElement(
			"tr",
			{ className: "" },
			React.createElement(
				"td",
				{ className: "" },
				React.createElement("input", {
					name: "name",
					id: "name",
					type: "text",
					className: "form-control",
					placeholder: this.props.accessGroup.name,
					readOnly: !this.state.editable,
					required: this.state.editable,
					value: this.state.name,
					onChange: function(e){
						_that.setState({name: e.target.value});
					}
				})
			),
			React.createElement(
				"td",
				{ className: "" },
				React.createElement("input", {
					name: "subnet",
					id: "subnet",
					type: "text",
					className: "form-control",
					placeholder: this.props.accessGroup.subnet,
					readOnly: !this.state.editable,
					required: this.state.editable,
					value: this.state.subnet,
					onChange: function(e){
						_that.setState({subnet: e.target.value});
					}
				})
			),
			React.createElement(
				"td",
				{ className: "" },
				React.createElement("input", {
					name: "mask",
					id: "mask",
					type: "text",
					className: "form-control",
					placeholder: this.props.accessGroup.mask,
					readOnly: !this.state.editable,
					required: this.state.editable,
					value: this.state.mask,
					onChange: function(e){
						_that.setState({mask: e.target.value});
					}
				})
			),
			React.createElement(
				"td",
				{ className: "" },
				React.createElement(
					"button",
					{ 
						className: "btn btn-info", 
						type: "button", 
						onClick: function(e) {
							_that.modify(e);
						}
					},
					this.state.editable ? "Speichern" : "Ändern"
				)
			),
			React.createElement(
				"td",
				{ className: "" },
				React.createElement(
					"button",
					{ 
						className: "btn btn-danger", 
						type: "button", 
						onClick: function(e) {
							_that.delete(e);
						}
					},
					"Löschen"
				)
			)
		)

	}
}
class AccessGroups extends React.Component {
	render() {
		return React.createElement(
			"table",
			{ className: "table table-striped" },
			React.createElement(
				"thead",
				null,
				React.createElement(
					"tr",
					null,
					React.createElement(
						"th",
						null,
						"Name"
					),
					React.createElement(
						"th",
						null,
						"subnet"
					),
					React.createElement(
						"th",
						null,
						"mask"
					)
				)
			),
			React.createElement(
				"tbody",
				null,
				this.props.accessGroups && this.props.accessGroups.map(function (accessGroup) {
					return React.createElement(
						AccessGroupLine,
						{ accessGroup: accessGroup, refresh: this.props.refresh }
					);
				}.bind(this))

			)
		);
	}

}

class App extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			accessGroups: [],
			subnet: '',
			id: '',
			mask: ''
		};

		this.create = this.create.bind(this);
		this.handleChange = this.handleChange.bind(this);
		this.refresh = this.refresh.bind(this);
	}
	refresh(){
		var _this = this
		// get all entities - GET
		fetch(window.location.toString() + "api/accessgroup", {
			"method": "GET",
			"headers": {
				//    "x-rapidapi-host": "fairestdb.p.rapidapi.com",
				//    "x-rapidapi-key": API_KEY
			}
		}).then(function (response) {
			return response.json();

		}).then(function (response) {
			_this.setState({
				accessGroups: response
			});
		}).catch(function (err) {
			console.log(err);
		});

	}
	componentDidMount() {
		this.refresh()
	}

	create(e) {
		// add entity - POST
		e.preventDefault();
		// creates entity
		fetch(window.location.toString() + "api/accessgroup", {
			"method": "POST",
			"headers": {
				"content-type": "application/json",
				"accept": "application/json"
			},
			"body": JSON.stringify({
				name: this.state.name,
				subnet: this.state.subnet,
				mask: this.state.mask
			})
		})
			.then(response => response.json())
			.then(response => {
				console.log(response)
			})
			.catch(err => {
				console.log(err);
			})
		.finally(() => this.refresh());
		;
	}

	handleChange(changeObject) {
		this.setState(changeObject)
	}

	render() {
		var _this4 = this;

		return React.createElement(
			"div",
			{ className: "container" },
			React.createElement(
				"div",
				{ className: "row justify-content-center" },
				React.createElement(
					"div",
					{ className: "col-md-8" },
					React.createElement(
						"h1",
						{ className: "display-4 text-center" },
						"Manage Access Groups"
					),
					React.createElement(
						"form",
						{ className: "d-flex flex-column" },
						React.createElement(
							"legend",
							{ className: "text-center" },
							"Add-Update-Delete AccessGroup"
						),
						React.createElement(
							"label",
							{ htmlFor: "name" },
							"AccessGroup Name:",
							React.createElement("input", {
								name: "name",
								id: "name",
								type: "text",
								className: "form-control",
								value: this.state.name,
								onChange: function onChange(e) {
									return _this4.handleChange({ name: e.target.value });
								},
								required: true
							})
						),
						React.createElement(
							"label",
							{ htmlFor: "subnet" },
							"AccessGroup subnet:",
							React.createElement("input", {
								name: "subnet",
								id: "subnet",
								type: "test",
								className: "form-control",
								value: this.state.subnet,
								onChange: function onChange(e) {
									return _this4.handleChange({ subnet: e.target.value });
								},
								required: true
							})
						),
						React.createElement(
							"label",
							{ htmlFor: "mask" },
							"AccessGroup subnet mask:",
							React.createElement("input", {
								name: "mask",
								id: "mask",
								type: "text",
								className: "form-control",
								value: this.state.mask,
								onChange: function onChange(e) {
									return _this4.handleChange({ mask: e.target.value });
								}
							})
						),
						React.createElement(
							"button",
							{ className: "btn btn-primary", type: "button", onClick: function onClick(e) {
								return _this4.create(e);
							} },
							"Add"
						)
					),
					React.createElement(AccessGroups, { accessGroups: this.state.accessGroups, refresh: this.refresh })
				)
			)
		);
	}
}

let domContainer = document.querySelector('#App');
ReactDOM.render(React.createElement(App, null), domContainer)
