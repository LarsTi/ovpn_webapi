class CertificateList extends React.Component{
	constructor(props){
		super(props);
		if(!props.apiUrl){
			console.log("Keine API URL gesetzt");
		}
	}
	refresh(){
		fetch(this.props.apiUrl).then( response => {
			if(response.ok){
				response.json().then(this.dataArrived.bind(this));
			}else{
				console.log("Fehler bei fetch!");
			}
		});
	}
	dataArrived(data){
		this.setState({data: data});
		console.log("received certificate data");
	}
	componentDidMount(){
		this.refresh();
	}
	create(){
		fetch(this.props.apiUrl, {
			"method": "POST",
			"headers": {
				"content-type": "application/json",
				"accept": "application/json"
			}
		}).then(response => {
			this.refresh()
		});
	}
	render(){
		return React.createElement("div", {},
			React.createElement("div", {className: "row"},
				React.createElement("div", {className: "col-2"}, 
					React.createElement("b", {}, "Common Name")),
				React.createElement("div", {className: "col-4"}, 
					React.createElement("b", {}, "Erzeugt am")),
				React.createElement("div", {className: "col-4"}, 
					React.createElement("b", {}, "Gültig bis")),
				React.createElement("div", {className: "col-2"}, 
					React.createElement("button", {
						onClick: this.create.bind(this)
					}, "Zertifikat anlegen")
				)
			),
			this.state && this.state.data && this.state.data.map( function(cert){
				cert.refresh = this.refresh.bind(this);
				cert.key = "cert-" + cert.ID;
				cert.apiUrl = this.props.apiUrl + "/" + cert.ID;
				return React.createElement(Certificate, cert);
			}.bind(this))
		)
	}
}
class Certificate extends React.Component{
	constructor(props){
		super(props);
		if(!props.apiUrl){
			console.log("Keine API URL verfügbar");
		}
	}
	revoke(){
		fetch(this.props.apiUrl, {
			"method": "DELETE",
			"headers": {
				"content-type": "application/json",
				"accept": "application/json"
			}
		}).then(response => {
			this.props.refresh()
		});
	}
	render(){
		return React.createElement(
			"div", {className: "row"},
			React.createElement(
				"div", {className: "col-2"},
				this.props.common_name),
			React.createElement(
				"div", {className: "col-4"},
				this.props.CreatedAt),
			React.createElement(
				"div", {className: "col-4"},
				this.props.valid_to),
			React.createElement(
				"div", {className: "col-1"},
				React.createElement("a", {href: this.props.apiUrl},
					React.createElement("button",{}, "Download")
				)
			),
			React.createElement(
				"div", {className: "col-1"},
				React.createElement("button", {
					onClick: this.revoke.bind(this)
				},"Revoke"))
		);
	}
}
