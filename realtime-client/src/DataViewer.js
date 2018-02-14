import React from 'react';
import Websocket from 'react-websocket';
import {
  LineChart
} from 'react-easy-chart';


class DataViewer extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      websocket_url: 'ws://localhost:8002/ws',
      componentWidth: 600,
      motion: [[{x:0,y:0}]],
      temperature: [[{x:0,y:0}]],
      ambient_light: [[{x:0,y:0}]],
      pressure: [[{x:0,y:0}]],
      humidity: [[{x:0,y:0}]],
    };
  }

  handleData(payload) {
    let result = JSON.parse(payload);
    let values = this.state[result.name]
    if (values[0].length > 60) {
      values[0].shift();
    }
    const largestX = values[0][values[0].length - 1].x
    values[0].push({
      x: largestX + 1,
      y: result.data
    })
    
    this.state[result.name] = values;
    this.setState(this.state)
  }

  render() {
    return ( 
    <div>
      <p> Motion </p> 
      <LineChart data = {this.state.motion}
      width = {this.state.componentWidth}
      height = {this.state.componentWidth / 2}
      axisLabels = {{x: 'Hour',y: 'Percentage'}}
      interpolate = {'cardinal'}
      // yDomainRange={[0, 100]}
      axes grid style = {{'.line0': {stroke: 'green'}}}
      /> 
      <p> Ambient Light </p> 
      <LineChart data = {this.state.ambient_light}
      width = {this.state.componentWidth}
      height = {this.state.componentWidth / 2}
      axisLabels = {{x: 'Hour',y: 'Percentage'}}
      interpolate = {'cardinal'}
      // yDomainRange={[0, 100]}
      axes grid style = {{'.line0': {stroke: 'green'}}} /> 
      <p> Temperature </p> 
      <LineChart data = {this.state.temperature}
      width = {this.state.componentWidth}
      height = {this.state.componentWidth / 2}
      axisLabels = {{x: 'Hour',y: 'Percentage'}}
      interpolate = {'cardinal'}
      // yDomainRange={[0, 100]}
      axes grid style = {{'.line0': {stroke: 'green'}}}
      /> 
      <p> Pressure </p> 
      <LineChart data = {this.state.pressure}
      width = {this.state.componentWidth}
      height = {this.state.componentWidth / 2}
      axisLabels = {{x: 'Hour',y: 'Percentage'}}
      interpolate = {'cardinal'}
      // yDomainRange={[0, 100]}
      axes grid style = {{'.line0': {stroke: 'green'}}}
      /> 
 
      <p> Humidity </p> 
      <LineChart data = {this.state.humidity}
      width = {this.state.componentWidth}
      height = {this.state.componentWidth / 2}
      axisLabels = {{x: 'Hour',y: 'Percentage'}}
      interpolate = {'cardinal'}
      // yDomainRange={[0, 100]}
      axes grid style = {{'.line0': {stroke: 'green'}}}
      /> 


      <Websocket url = {this.state.websocket_url}      onMessage = {this.handleData.bind(this)}
      /> 
    </div>
    );
  }
}

export default DataViewer;