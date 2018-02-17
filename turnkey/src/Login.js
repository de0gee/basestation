import React, { Component } from "react";
import { Button, FormGroup, FormControl, ControlLabel } from "react-bootstrap";
import ToggleDisplay from 'react-toggle-display';
import "./Login.css";
import axios from 'axios'

export default class Login extends Component {
  constructor(props) {
    super(props);

    this.state = {
      email: "",
      ssid: "",
      password: "",
      show: true
    };
  }

  validateForm() {
    return this.state.ssid.length > 0 && this.state.email.length > 0 && this.state.password.length > 7;
  }

  handleChange = event => {
    this.setState({
      [event.target.id]: event.target.value
    });
  }

  handleSubmit = event => {
    event.preventDefault();
    console.log(this.state);
    var currentState = this.state;
    axios.post('/signin', {
      currentState
    })
    .then(function (response) {
      console.log(response);
    })
    .catch(function (error) {
      console.log(error);
    });
  
    this.setState({
      show: false
    })

  }

  render() {
    return (
      <div className="Login">
        <ToggleDisplay show={this.state.show}>
          <form onSubmit={this.handleSubmit}>
          <FormGroup controlId="email" bsSize="large">
              <ControlLabel>Email</ControlLabel>
              <FormControl
                autoFocus
                value={this.state.email}
                onChange={this.handleChange}
                type="email"
              />
            </FormGroup>
            <FormGroup controlId="ssid" bsSize="large">
              <ControlLabel>SSID</ControlLabel>
              <FormControl
                autoFocus
                value={this.state.ssid}
                onChange={this.handleChange}
              />
            </FormGroup>
            <FormGroup controlId="password" bsSize="large">
              <ControlLabel>Password</ControlLabel>
              <FormControl
                value={this.state.password}
                onChange={this.handleChange}
              />
            </FormGroup>
            <Button
              block
              bsSize="large"
              disabled={!this.validateForm()}
              type="submit"
            >
              Login
            </Button>
          </form>
        </ToggleDisplay>
        <ToggleDisplay if={!this.state.show}>
        <div className="jumbotron">
          <h1>Thanks.</h1> 
          <p>Please wait about 2 minutes and check your email about the login.</p> 
        </div>
          </ToggleDisplay>
      </div>
    );
  }
}