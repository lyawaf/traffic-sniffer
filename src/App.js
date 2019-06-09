import React, {Component} from 'react';
import './App.css';
import Sessions from './Sessions'
class App extends Component {

	state = {
		data: []
	}

	componentDidMount() {
		fetch('http://localhost:9999/',{'method':'POST'})
		.then((response)=>{response.json().then((data)=> {
			//let temp = data
			this.setState({data:data})})})
	}


	


  	render() {
  	//console.log("state: ", this.state.data)
    return(
          <div>
          	<Sessions respdata={this.state.data}/>
          </div>
          )}
  
}

export default App;
