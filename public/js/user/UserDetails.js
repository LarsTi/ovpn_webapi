class UserDetails extends React.Component{
	constructor(props){
		super(props);
		this.state = {
			active: false,
			deleteLabel: "Benutzer komplett löschen",
			deleteAllowed: false,
			data: {
				name: props.name,
				surname: props.surname,
				passwd: props.passwd,
				org: props.org,
				mail: props.mail,
			}
		}
	}
	getNameElement(name){
		return React.createElement(
			"input", {
				type: "text",
				placeholder: name,
				value: this.state.data[name],
				onChange: function(e){
					var d = this.state.data;
					d[name] = e.target.value;
					this.setState({data:d});

					if( d.name === this.props.name
						&& d.surname === this.props.surname
						&& d.passwd === this.props.passwd
						&& d.org === this.props.org
						&& d.mail === this.props.mail ){
						this.setState({active: false});
					}else{
						this.setState({active: true});
					}
				}.bind(this)
			}
		);
	}
	onSave(){
		if( this.state.active){
			fetch(this.props.apiUrl + "/" + this.props.ID, {
				"method": "PUT",
				"headers": {
					"content-type": "application/json",
					"accept": "application/json"
				},
				"body": JSON.stringify(this.state.data)
			}).then(this.props.app.replaceView(this.props.app.state.user));
		}
	}
	onDelete(){
		if( this.state.deleteAllowed){
			fetch(this.props.apiUrl + "/" + this.props.ID, {
				"method": "DELETE",
				"headers": {
					"content-type": "application/json",
					"accept": "application/json"
				}
			}).then(this.props.app.replaceView(this.props.app.state.user));
		}else{
			this.setState({
				deleteAllowed: true,
				deleteLabel: "Wirklich löschen"
			})
		}
	}
	getChangeLine(){
		return React.createElement("div", {className: "row"},
			React.createElement("div", {className: "col-2"},
				this.getNameElement("surname")),
			React.createElement("div", {className: "col-2"},
				this.getNameElement("name")),
			React.createElement("div", {className: "col-2"},
				this.getNameElement("passwd")),
			React.createElement("div", {className: "col-1"},
				this.getNameElement("org")),
			React.createElement("div", {className: "col-3"},
				this.getNameElement("mail")),
			React.createElement("div", {className: "col-2"},
				React.createElement("button",{
					onClick: this.onSave.bind(this),
					disabled: !this.state.active
				}, "Save"))
		);
	}
	getDeleteLine(){
		return React.createElement("div", {className: "row"},
			React.createElement("div", {className: "col-12"},
				React.createElement("button",{
					onClick: this.onDelete.bind(this)
				}, this.state.deleteLabel))
		);
	}
	getCertificateListHeader(){
		return React.createElement("div", {className: "row"},
			React.createElement("div", {className: "col-12"},
				React.createElement("h3", {}, "Zertifikate:"))
		);
	}

	render(){
		return React.createElement("div", {},
			this.getChangeLine(),
			this.getDeleteLine(),
			this.props.app.getSeperator(),
			React.createElement(CertificateList, {
				apiUrl: this.props.apiUrl + "/" + this.props.ID + "/certificate",
				data: this.props
			}),
			this.props.app.getSeperator(),
			React.createElement(AccessList, {
				app: this.props.app,
				apiUrl: this.props.apiUrl + "/" + this.props.ID + "/access",
				user: this.props.ID,
				apiAccess: "/api/accessgroup"
			})
		);

	}
}
