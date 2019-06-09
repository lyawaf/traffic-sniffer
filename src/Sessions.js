import React, {Component} from 'react';
import Session from './Session'

class Sessions extends Component {
  render() {
  	let data = this.props.respdata
  	let sessionsList
  	//console.log('data: ', data)
  	if(data.length) {
  		//console.log(data[0]._id)
		sessionsList = data.map((index)=>
			<div key={index._id}>
				<Session data={index}/>
			</div>
		)}
	else
		sessionsList = <div>Loading...</div>
    return(
      <div>
      {sessionsList}
      </div>
      )
  }
}

export default Sessions;
