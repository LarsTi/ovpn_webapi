class Access extends React.Component{
	constructor(props){
		super(props);
		if(!props.apiUrl){
			console.log("Keine API URL gesetzt");
		}
	}
	buttonClick(){
		if(this.props.allowed){
			//zugriff ist momentan erlaubt, also revoken
			fetch(this.props.apiUrl + "/" + this.props.id, {
				"method": "DELETE",
				"headers": {
					"content-type": "application/json",
					"accept": "application/json"
				}
			}).then(response => {
				this.props.refresh();
			});
		}else{
			//Zugriff ist momentan nicht erlaubt, also erlauben
			fetch(this.props.apiUrl, {
				"method": "POST",
				"headers": {
					"content-type": "application/json",
					"accept": "application/json"
				},
				"body": JSON.stringify({
					User: this.props.user,
					Group: this.props.id
				})
			}).then(response => {
				this.props.refresh();
			});
		}
	}
	render(){
		return React.createElement("div", {className: "row"},
			React.createElement("div", {className:"col-6"}, this.props.name),
			React.createElement("div", {className:"col-4"}, this.props.subnet + "/" + this.props.mask),
			React.createElement("div", {className: "col-2"},
				React.createElement("button", {
					onClick: this.buttonClick.bind(this)
				}, this.props.text)
			)
		)
	}
}
