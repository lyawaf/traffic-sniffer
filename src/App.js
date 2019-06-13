import React, {Component} from 'react';
import './App.css';
import Sessions from './Sessions'
import AddLabelForm from './AddLabelForm'
import Labels from './Labels'
class App extends Component {

	state = {
		data: [],
		currentLabel: {},
		lastUpdate : 0
	}

	componentDidMount() {
		let data = JSON.stringify({"lastUpdate": 0})
		fetch('http://localhost:9999/',{'method':'POST','body':data})
		.then((response)=>{response.json().then((data)=> {
			//let temp = data
			//console.log("data= ", data[data.length-1].last_update)
			this.setState({data:data})})})
	}

	componentDidUpdate(prevState) {
	let data = JSON.stringify({"lastUpdate": 0})
  if(prevState.data!==this.state.data) {
    fetch('http://localhost:9999',{'method':'POST','body':data})
    .then((response)=>{response.json().then((data)=> {
      //let temp = data
       
      //console.log("data= ", data)
      this.setState({data:data})}


      )})
  }
}

	chngLabel = (e) => {
		this.setState({currentLabel:e})
	}


	


  	render() {
  //	console.log("state: ", this.state.lastUpdate)
    return(
          <div>
          	<Labels/>
          	<Sessions respdata={this.state.data}/>
          	<AddLabelForm/>
          	
          </div>
          )}
  
}

export default App;
