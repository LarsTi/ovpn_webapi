class AccessList extends React.Component{
	constructor(props){
		super(props)
		if(!props.apiUrl){
			console.log("Keine API URL gesetzt");
		}
	}
	getHeader(){
		return React.createElement("div", {className: "row"},
			React.createElement("div", {className: "col-6"},
				React.createElement("b", {}, "Name")),
			React.createElement("div", {className: "col-4"},
				React.createElement("b", {}, "Subnet/Mask")),
			React.createElement("div", {className: "col-2"})
		);
	}
	completeRefresh(){
		fetch(this.props.apiAccess).then( response => {
			if(response.ok){
				response.json().then(this.allData.bind(this));
			}else{
				console.log("Fehler bei fetch!");
			}
		});
		fetch(this.props.apiUrl).then( response =>{
			if(response.ok){
				response.json().then(this.userData.bind(this));
			}else{
				console.log("Fehler bei fetch!");
			}
		});
	}
	userData(data){
		if( data && data.length > 0 && this.state && this.state.data ){
			var arr = this.state.data;
			var changed = false;
			for(var i = 0; i < data.length; i++){
				for(var j = 0; j < this.state.data.length; j++){
					if( arr[j].id == data[i].group){
						changed = true;
						arr[j].allowed = true;
						arr[j].text = "Zugriff entziehen";
					}
				}
			}
			if(changed){
				this.setState({data: arr});
			}
		}
	}
	allData(data){
		if(data && data.length > 0){
			var arr = [];
			for(var i = 0; i < data.length; i++){
				var d = {
					name: data[i].name,
					id: data[i].ID,
					subnet: data[i].subnet,
					mask: data[i].mask,
					allowed: false,
					text: "Zugriff erlauben"
				}
				arr.push(d)
			}
			this.setState({data: arr});
		}
	}
	componentDidMount(){
		this.completeRefresh();
	}
	render(){
		return React.createElement("div", {}, 
			this.getHeader(),
			this.props.app.getSeperator(),
			this.state && this.state.data && this.state.data.map( function (ag){
				ag.key = "usergroup-" + ag.id;
				ag.apiUrl = this.props.apiUrl;
				ag.refresh = this.completeRefresh.bind(this);
				ag.user = this.props.user
				return React.createElement(Access, ag);
			}.bind(this))
		);
	}
}
