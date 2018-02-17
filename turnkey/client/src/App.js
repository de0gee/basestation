import React, { Component } from 'react';
import './App.css';
import Login from './Login.js';

class App extends Component {
  render() {
    return (
      <div className="App">
        <div className="jumbotron">
          <h1>Welcome</h1>
          <p>Please enter your credentials</p>
        </div>
        <Login />
      </div>
    );
  }
}

export default App;
