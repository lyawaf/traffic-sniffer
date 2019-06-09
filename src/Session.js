import React, {Component} from 'react';
import { Base64 } from 'js-base64';

class Session extends Component {
  state = {
  	isOpen: false
  }

  click = () =>{
  	this.setState({isOpen:!this.state.isOpen})
  }

  render() {

  	let data = this.props.data
  	let sessionList = data.session.packets.map((index)=>
  		<p className={'session '+index.owner}>{Base64.decode(index.data)}</p>)
  	
  	let info = <div>{sessionList}</div>
  	if(this.state.isOpen) {
    return(
      <div>
      <button onClick={this.click}>{data._id}</button>
      {info}
      </div>
      )}
    else {
    	return(
      <div>
      <button onClick={this.click}>{data._id}</button>
      </div>
      )
    }
  
}
}

export default Session;
