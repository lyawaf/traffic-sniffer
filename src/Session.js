import React, {Component} from 'react';
import { Base64 } from 'js-base64';

class Session extends Component {
  state = {
  	isOpen: false
  }


  hd = (buffer, blockSize)=> {
    blockSize = blockSize || 16;
    var lines = [];
    var hex = "0123456789ABCDEF";
    for (var b = 0; b < buffer.length; b += blockSize) {
        var block = buffer.slice(b, Math.min(b + blockSize, buffer.length));
        var addr = ("0000" + b.toString(16)).slice(-4);
        var codes = block.split('').map(function (ch) {
            var code = ch.charCodeAt(0);
            return " " + hex[(0xF0 & code) >> 4] + hex[0x0F & code];
        }).join("");
        codes += " 00".repeat(blockSize - block.length)+"\t";
        var chars = block.replace(/[\x00-\x1F\x20]/g, '.');
        chars +=  " ".repeat(blockSize - block.length);
        lines.push(addr + "\t" + codes + "\t" + chars);
    }
    return lines.join("\n");
}

  click = () =>{
  	this.setState({isOpen:!this.state.isOpen})
  }

  createSessionList = (dataArr) => {
    return(
      dataArr.session.packets.map((element, index)=>
      <p key={index} className={'session '+element.owner}>{this.hd(Base64.decode(element.data).replace(/[^\x20-\x7E]/g, "."),32)}</p>)
      )
  }

  render() {

  	let data = this.props.data
  	let port = data.session.clientport
    let addr = data.session.clientaddr
  	let info = <div>{this.createSessionList(data)}</div>
  	if(this.state.isOpen) {
    return(
      <div className="sess-cont">
      
      <p className="label-logo">label name </p>
      <div className="sess-info">
      <button className="session-button" onClick={this.click}>Скрыть сессию</button>
      <p className="client-info">Port: {port} <br/> Client addr: {addr}</p>
     </div>
      {info}
      </div>
      )}
    else {
    	return(
      <div className="sess-cont">
      
      <p className="label-logo">label name </p>
      <div className="sess-info">
      <button className="session-button" onClick={this.click}>Показать сессию</button>
      <p className="client-info">Port: {port} <br/> Client addr: {addr}</p>
      </div>
      </div>
      )
    }
  
}
}

export default Session;
