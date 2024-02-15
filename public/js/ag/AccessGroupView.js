class AccessGroupView extends React.Component{
	constructor(props){
		super(props)
		if(!props.apiUrl){
			console.log("Keine API URL gesetzt");
		}
	}
	refresh(){
		fetch(this.props.apiUrl).then( response => {
			if(response.ok){
				response.json().then(this.dataArrived.bind(this));
			}else{
				console.log("Fehler bei fetch!")
			}
		});
	}
	dataArrived(data){
		this.setState({data: data});
		console.log("received AccessGroup data");
	}
	componentDidMount() {
		this.refresh();
	}
	getHeader(){
		return React.createElement(
			"div", {className: "row"},
			React.createElement(
				"div",{className: "col-4"},
				React.createElement("b", {}, "Bezeichnung")
			),
			React.createElement(
				"div",{className: "col-4"},
				React.createElement("b", {}, "Subnet")
			),
			React.createElement(
				"div",{className: "col-2"},
				React.createElement("b", {}, "Mask")
			),
			React.createElement(
				"div",{className: "col-2"}
			)
		
		);
	}
	render() {
		return React.createElement(
			"div", {className: "container"}, 
			this.props.app.getTitle("Access Groups pflegen"),
			this.props.app.getSeperator(),
			this.getHeader(),
			this.props.app.getSeperator(),
			React.createElement(AccessGroupNew, {
				apiUrl: this.props.apiUrl,
				refresh: this.refresh.bind(this)
			}),
			this.props.app.getSeperator(),
			this.state && this.state.data && this.state.data.map( function(ag){
				ag.refresh = this.refresh.bind(this);
				ag.key = "accessgroup-" + ag.ID;
				ag.apiUrl = this.props.apiUrl;
				return React.createElement(
					AccessGroup,
					ag,
				);
			}.bind(this))
		);
	}

}
