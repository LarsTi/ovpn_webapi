'use strict';

class UserLine extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			name: "",
			surname: "",
			org: "",
			mail: "",
			passwd: '',
			editable: false
		};
		this.modify = this.modify.bind(this);
		this.delete = this.delete.bind(this);
	}
	modify(e) {
		e.preventDefault();
		if (this.state.editable){
			if ( this.state.name === this.props.user.name
				&& this.state.surname === this.props.user.surname
				&& this.state.org === this.props.user.org
				&& this.state.mail === this.props.user.mail
				&& this.state.passwd === this.props.user.passwd ){
				//Es gab keine Änderung.
				this.setState({editable: false });
			}else{
				//Updaten der Änderung ins Backend
				// this will update entries with PUT
				fetch(window.location.protocol + "//" + window.location.host + "/api/user/" + this.props.user.ID, {
					"method": "PUT",
					"headers": {
						"content-type": "application/json",
						"accept": "application/json"
					},
					"body": JSON.stringify({
						name: this.state.name,
						surname: this.state.surname,
						name: this.state.name,
						org: this.state.org,
						mail: this.state.mail,
						passwd: this.state.passwd,
						ID: this.props.user.ID
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
				name:this.props.user.name,
				surname: this.props.user.surname,
				org: this.props.user.org,
				mail: this.props.user.mail,
				passwd: this.props.user.passwd,
				editable: true
			})
		}
	}
	delete(e) {
		// delete entity - DELETE
		e.preventDefault();
		// deletes entities
		fetch(window.location.protocol + "//" + window.location.host + "/api/user/" + this.props.user.ID, {
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
					placeholder: this.props.user.name,
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
					name: "surname",
					id: "surname",
					type: "text",
					className: "form-control",
					placeholder: this.props.user.surname,
					readOnly: !this.state.editable,
					required: this.state.editable,
					value: this.state.surname,
					onChange: function(e){
						_that.setState({surname: e.target.value});
					}
				})
			),
			React.createElement(
				"td",
				{ className: "" },
				React.createElement("input", {
					name: "org",
					id: "org",
					type: "text",
					className: "form-control",
					placeholder: this.props.user.org,
					readOnly: !this.state.editable,
					required: this.state.editable,
					value: this.state.org,
					onChange: function(e){
						_that.setState({org: e.target.value});
					}
				})
			),
			React.createElement(
				"td",
				{ className: "" },
				React.createElement("input", {
					name: "mail",
					id: "mail",
					type: "text",
					className: "form-control",
					placeholder: this.props.user.mail,
					readOnly: !this.state.editable,
					required: this.state.editable,
					value: this.state.mail,
					onChange: function(e){
						_that.setState({mail: e.target.value});
					}
				})
			),
			React.createElement(
				"td",
				{ className: "" },
				React.createElement("input", {
					name: "passwd",
					id: "passwd",
					type: "text",
					className: "form-control",
					placeholder: this.props.user.passwd,
					readOnly: !this.state.editable,
					required: this.state.editable,
					value: this.state.passwd,
					onChange: function(e){
						_that.setState({passwd: e.target.value});
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
class User extends React.Component {
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
						"Surname"
					),
					React.createElement(
						"th",
						null,
						"Organisation"
					),
					React.createElement(
						"th",
						null,
						"Mail"
					),
					React.createElement(
						"th",
						null,
						"Passwd (hash)"
					)
				)
			),
			React.createElement(
				"tbody",
				null,
				this.props.user && this.props.user.map(function (u) {
					return React.createElement(
						UserLine,
						{ 
							user: u, 
							refresh: this.props.refresh,
							key: u.ID
						}
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
			user: [],
			name: "",
			surname: "",
			org: "",
			mail: "",
			passwd: ""
		};

		this.create = this.create.bind(this);
		this.handleChange = this.handleChange.bind(this);
		this.refresh = this.refresh.bind(this);
	}
	refresh(){
		var _this = this
		// get all entities - GET
		fetch(window.location.protocol + "//" + window.location.host + "/api/user", {
			"method": "GET",
			"headers": {
				//    "x-rapidapi-host": "fairestdb.p.rapidapi.com",
				//    "x-rapidapi-key": API_KEY
			}
		}).then(function (response) {
			return response.json();

		}).then(function (response) {
			_this.setState({
				user: response
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
		fetch(window.location.protocol + "//" + window.location.host + "/api/user", {
			"method": "POST",
			"headers": {
				"content-type": "application/json",
				"accept": "application/json"
			},
			"body": JSON.stringify({
				name: this.state.name,
				surname: this.state.surname,
				org: this.state.org,
				mail: this.state.mail,
				passwd: this.state.passwd
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
		var _this = this;

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
						"Manage User"
					),
					React.createElement(
						"form",
						{ className: "d-flex flex-column" },
						React.createElement(
							"legend",
							{ className: "text-center" },
							"Add-Update-Delete User"
						),
						React.createElement(
							"label",
							{ htmlFor: "name" },
							"Name:",
							React.createElement("input", {
								name: "name",
								id: "name",
								type: "text",
								className: "form-control",
								value: this.state.name,
								onChange: function onChange(e) {
									return _this.handleChange({ name: e.target.value });
								},
								required: true
							})
						),
						React.createElement(
							"label",
							{ htmlFor: "surname" },
							"Surname:",
							React.createElement("input", {
								name: "surname",
								id: "surname",
								type: "text",
								className: "form-control",
								value: this.state.surname,
								onChange: function onChange(e) {
									return _this.handleChange({ surname: e.target.value });
								},
								required: true
							})
						),
						React.createElement(
							"label",
							{ htmlFor: "org" },
							"User Organisation:",
							React.createElement("input", {
								name: "org",
								id: "org",
								type: "text",
								className: "form-control",
								value: this.state.org,
								onChange: function onChange(e) {
									return _this.handleChange({ org: e.target.value });
								}
							})
						),
						React.createElement(
							"label",
							{ htmlFor: "mail" },
							"User Mail-address:",
							React.createElement("input", {
								name: "mail",
								id: "mail",
								type: "text",
								className: "form-control",
								value: this.state.mail,
								onChange: function onChange(e) {
									return _this.handleChange({ mail: e.target.value });
								}
							})
						),
						React.createElement(
							"label",
							{ htmlFor: "passwd" },
							"User Passwd (hashed):",
							React.createElement("input", {
								name: "passwd",
								id: "passwd",
								type: "text",
								className: "form-control",
								value: this.state.passwd,
								onChange: function onChange(e) {
									return _this.handleChange({ passwd: e.target.value });
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
					React.createElement(User, 
						{ 
							user: this.state.user, 
							refresh: this.refresh 
						}
					)
				)
			)
		);
	}
}

let domContainer = document.querySelector('#App');
ReactDOM.render(React.createElement(App, null), domContainer)
