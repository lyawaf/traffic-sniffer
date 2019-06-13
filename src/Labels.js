import React, {Component} from 'react';
import Label from './Label'


class Labels extends Component {


state = {
  labels:[]
}


componentDidMount() {
    fetch('http://localhost:9999/getLabels')
    .then((response)=>{response.json().then((data)=> {
      //let temp = data
      this.setState({labels:data.L})})})
  }

componentDidUpdate(prevState) {
  if(prevState.labels!==this.state.labels) {
    fetch('http://localhost:9999/getLabels')
    .then((response)=>{response.json().then((data)=> {
      //let temp = data
      this.setState({labels:data.L})})})
  }
}


  render() {
    console.log("from labels")
  	let stateLabels = this.state.labels
    let labels
    if(stateLabels.length) {
      labels = stateLabels.map((index)=>
        <div key={index.name} className="labels">
        <Label data={index}/>
        </div>
        )
    }
    else {
      labels = <div></div>
    }

    return(
      <div className="label-cont">{labels}</div>
      )
  }
}

export default Labels;
