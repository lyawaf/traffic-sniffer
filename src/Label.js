import React, {Component} from 'react';
import { Base64 } from 'js-base64';


class Label extends Component {
  render() {
  	let data = this.props.data
  	let name = data.Name
  	let type = data.Type
  	let color = data.Color
  	let regexp = Base64.decode(data.RawRegexp)
  	var style
  	if(color==="#ffffff"||color==="white") {
  	style ={ backgroundColor:color,
  				   'color':'black'}
  	}
  	else {
  	style =	{ backgroundColor:color,
  				   'color':'white'}	
  	}

  	var shortRegexp = regexp.length<18?regexp:regexp.slice(0,18)

  	return(
	  	<div className={regexp.length<18?"holder":"holder holder-with-block"}>
	  		<button style={style} className="label">{name}<br/>{type}<br/>{shortRegexp}</button>
	  		<div className="hide-block">{regexp}</div>
	  	</div>
	)
  }
}

export default Label;
