"use-strict";
class App extends React.Component{
	constructor(props){
		super(props);
		
		this.state = {
			activeView: "user",
			accessGroup: React.createElement(AccessGroupView, {
				apiUrl: "/api/accessgroup",
				app: this
			}),
			user: React.createElement(UserView, {
				apiUrl: "/api/user",
				app: this
			})
		}
		
	}
	replaceView(element, functionToCall){
		if( !element ){
			console.log("Keine Änderung am Element vornehmen");
		}else if( element === this.state.user){
		//Force state update
			this.setState({rerender: !this.state.rerender});
			this.setState({activeView: "user"})
		}else if(this.state.activeView === "customView"){
		//Ändern des view namens, damit definitiv änderung des states existiert
			
			this.setState({
				activeView: "customView2",
				customView: undefined,
				customView2: element
			});
		}else{
			this.setState({
				activeView: "customView",
				customView: element,
				customView2: undefined
			});
		}
		if(functionToCall){
			functionToCall();
		}
	}
	getHeader(){
		return React.createElement("div", {className: "row"},
			React.createElement("div", {className: "col-4"},
				React.createElement("button",{
					onClick: function(e){
						this.setState({activeView: "accessGroup"});
					}.bind(this)
				},"Access Groups")
			),
			React.createElement("div", {className: "col-4"},
				React.createElement("button",{
					onClick: function(e){
						this.setState({activeView: "user"});
					}.bind(this)
				},"User")
			)
		);
	}
	getTitle(title){
		return React.createElement(
			"div", {className: "row"},
			React.createElement("div", {className: "col-12"},
				React.createElement("h2", null, title)
			)
		);
	}
	getSeperator(){
		return React.createElement(
			"div", {className: "row"},
			React.createElement("div", {className: "col-12"},
				React.createElement("hr", null)
			)
		);
	}
	render(){
		return React.createElement(
			"div", {className: "container"},
				this.getHeader(),
				this.state[this.state.activeView]
			);
	}
}

ReactDOM.render(React.createElement(App, {}), document.getElementById("App"));
