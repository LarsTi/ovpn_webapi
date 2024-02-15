class UserNew extends React.Component{
	constructor(props){
		super(props)
		this.state = {
			active: false,
			data: {
				name: "",
				surname: "",
				org: "",
				mail: "",
				passwd: ""
			}
		}
		if(!props.apiUrl){
			console.log("Keine API URL gesetzt");
		}
	}
	onSave(e){
		if(this.state.active){
			fetch(this.props.apiUrl, {
				"method": "POST",
				"headers": {
					"content-type": "application/json",
					"accept": "application/json"
				},
				"body": JSON.stringify(this.state.data)
			}).then(this.props.app.replaceView(this.props.app.state.user, this.props.refresh));

			this.setState({data: {
				name: "",
				surname: "",
				org: "",
				mail: "",
				passwd: ""
			}});
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

					if(d.name !== "" && d.subnet !== "" && d.mask !== ""){
						this.setState({active: true});
					}else{
						this.setState({active: false});
					}
				}.bind(this)
			}
		)
	}

	render(){
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
				React.createElement("button", {
					onClick: this.onSave.bind(this),
					disabled: !this.state.active
				}, "Save")
			)
		);
	}
}
