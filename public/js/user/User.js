class User extends React.Component{
	constructor(props){
		super(props);
		if(!props.app.replaceView){
			console.log("Kein replaceView übergeben");
		}
	}
	onDetails(){
		this.props.app.replaceView(React.createElement(UserDetails, this.props));
	}
	render() {
		return React.createElement(
			"div", {className: "row flat-cols"},
			React.createElement(
				"div", {className: "col-12"},
				React.createElement("hr")
			),
			React.createElement(
				"div", {className: "col-2"},
				this.props.surname,
			),
			React.createElement(
				"div", {className: "col-2"},
				this.props.name
			),
			React.createElement(
				"div", {className: "col-2 overflow-auto"},
				this.props.passwd
			),
			React.createElement(
				"div", {className: "col-1"},
				this.props.org
			),
			React.createElement(
				"div", {className: "col-3"},
				this.props.mail
			),
			React.createElement(
				"div", {className: "col-2"},
				React.createElement("button", {onClick: this.onDetails.bind(this)}, "Details")
			),
		)

	}
}
