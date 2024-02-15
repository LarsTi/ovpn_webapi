class AccessGroup extends React.Component{
	constructor(props){
		super(props);
		this.state = {
			editable: false,
			data: {
				ID: this.props.ID,
				name: this.props.name,
				subnet: this.props.subnet,
				mask: this.props.mask
			},
			editButtonText: "Edit"
		}
		if(!props.apiUrl){
			console.log("Keine API URL gesetzt");
		}
	}
	onEdit(e){
		if(!this.props.ID){
			console.log("props.ID not found!");
			return;
		}
		if(this.state.editable){
			this.setState({editButtonText: "Edit"});
			fetch(this.props.apiUrl + "/" + this.props.ID, {
				"method": "PUT",
				"headers": {
					"content-type": "application/json",
					"accept": "application/json"
				},
				"body": JSON.stringify(this.state.data)
			}).then(this.props.refresh());
		}else{
			this.setState({editButtonText: "Save"});
		}

		this.setState({
			editable: !this.state.editable
		})
	}
	onDelete(e){
		if(!this.props.ID){
			console.log("props.ID not found!");
			return;
		}
		fetch(this.props.apiUrl + "/" + this.props.ID, {
			"method": "DELETE",
		}).then(this.props.refresh());
	}
	getNameElement(name){
		if(this.state.editable){
			return React.createElement(
				"input", {
					type: "text", 
					placeholder: this.props[name],
					readOnly: !this.state.editable,
					value: this.state.data[name],
					onChange: function(e){
						var d = this.state.data;
						d[name] = e.target.value;
						this.setState({data: d});
					}.bind(this)
				}
			)
		}else{
			return this.props[name]
		}
	}
	render() {

		return React.createElement(
			"div", {className: "row flat-cols"},
			React.createElement(
				"div", {className: "col-12"},
				React.createElement("hr")
			),
			React.createElement(
				"div", {className: "col-4"},
				this.getNameElement("name"),
			),
			React.createElement(
				"div", {className: "col-4"},
				this.getNameElement("subnet")
			),
			React.createElement(
				"div", {className: "col-2"},
				this.getNameElement("mask")
			),
			React.createElement(
				"div", {className: "col-1"},
				React.createElement("button", {onClick: this.onEdit.bind(this)}, this.state.editButtonText)
			),
			React.createElement(
				"div", {className: "col-1"},
				React.createElement("button", {onClick: this.onDelete.bind(this)}, "Delete")
			),
		)
	}
}
