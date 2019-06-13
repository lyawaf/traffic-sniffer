import React, {Component} from 'react';
import { Base64 } from 'js-base64';

class AddLabelForm extends Component {
  state = {
    name: '',
    type: 'in',
    regexp: '',
    color: '',
 } 


 onChangeName = (e) => {
  this.setState({name:e.target.value})
 }

 onChangeType = (e) => {
  this.setState({type:e.target.value})
 }

 onChangeRegexp = (e) => {
  this.setState({regexp:e.target.value})
 }

 onChangeColor = (e) => {
  this.setState({color:e.target.value})
 }


onSubmit = (e) => {
  e.preventDefault()
  //console.log(e.target[0])
  this.addLabel()
  //e.target.value=''
  for(var i=0;i<4;i++){
    if(i!==1)
      e.target[i].value = ''
  }
 // this.setState({name:'',type:'in',regexp:'',color:''})
  console.log(this.state)
} 

addLabel = () => {
  let reg = Base64.encode(this.state.regexp)
  let data = JSON.stringify({"name":this.state.name, "color":this.state.color , "type":this.state.type, "regexp":reg})
  fetch('http://localhost:9999/addLabel',{'method':'POST', 'body':data})
    .then((response)=>{response.text().then((data)=> {
      //let temp = data
      console.log(data)})})
}


  render() {
    
    return(
      <form onSubmit={this.onSubmit} className="form">
          <p><label> Name: <input type="text" name="name" value={this.state.name}
                           onChange={this.onChangeName}/></label></p>
          <p><label> Type: <select className="select" name="type"  value={this.state.type} onChange={this.onChangeType}>
              <option value="in">In</option>
              <option value="out">Out</option>
            </select>
              </label>
          </p>
          <p><label> Regexp: <input type="text" name="regexp" value={this.state.regexp}
                            onChange={this.onChangeRegexp}/></label></p>                 
          <p><label> Color: <input type="text" name="color" value={this.state.color}
                            onChange={this.onChangeColor}/></label></p>                  
          <p><input className="input" type="submit" value="Отправить" /></p>
      </form>
      )


  	}
}

export default AddLabelForm;
