import * as React from 'react';
import { RadioGroup, RadioButton } from 'react-radio-buttons'; 
import './ChooseActivity.css'
import axios from 'axios';
import Alert from 'react-s-alert';
import 'react-s-alert/dist/s-alert-default.css';
import 'react-s-alert/dist/s-alert-css-effects/flip.css';

class ChooseActivity extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
          value: "none",
       };
      var self = this;
      console.log(window.location.href.replace(":3000",":8002")+'activity');
        axios.get(window.location.href.replace(":3000",":8002")+'activity')
        .then(function (response) {
          console.log(response);
          if (response.data.success == true) {
            console.log(response.data.message);
            self.setState({value: response.data.message});
          } else {
            Alert.warn(response.data.message, {
              position: 'top-right',
              effect: 'flip',
              timeout: 3000
            });
          }
        })
        .catch(function (error) {
          console.log(error);
          Alert.error(error, {
              position: 'top-right',
              effect: 'flip',
              timeout: 3000
          });
        });


      }

    onChange(value) {
      console.log(value);
      axios.post(window.location.href.replace(":3000",":8002")+'activity', {
        activity: value
      })
      .then(function (response) {
        console.log(response);
        if (response.data.success == true) {
          Alert.success(response.data.message, {
            position: 'top-right',
            effect: 'flip',
            timeout: 3000
          });
        } else {
          Alert.warn(response.data.message, {
            position: 'top-right',
            effect: 'flip',
            timeout: 3000
          });
        }
      })
      .catch(function (error) {
        console.log(error);
        Alert.error(error, {
            position: 'top-right',
            effect: 'flip',
            timeout: 3000
        });
      });
    }
    
    render() {
        return (
          <div style={ { padding: 16 } } className="activities">
          <Alert stack={{limit: 3}} html={true} />
            <h4 style={ { marginTop: 32 } }>Classify Activity</h4>

            <RadioGroup onChange={ this.onChange } value={this.state.value}>
          <RadioButton value="none" pointColor="#999999">
            none
          </RadioButton>
          <RadioButton value="walking">
          walking
          </RadioButton>
          <RadioButton value="running">
          running
          </RadioButton>
          <RadioButton value="eating">
          eating
          </RadioButton>
          <RadioButton value="playing">
          playing
          </RadioButton>
          <RadioButton value="sleeping">
          sleeping
          </RadioButton>
          <RadioButton value="barking">
          barking
          </RadioButton>
        </RadioGroup>
          </div>
        );
      }
}

export default ChooseActivity;