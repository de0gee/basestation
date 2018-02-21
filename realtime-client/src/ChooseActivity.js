import * as React from 'react';
import { RadioGroup, RadioButton } from 'react-radio-buttons'; 

class ChooseActivity extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
           value: "none",
        };
      }

      onChange(value) {
        console.log(value);
      }
    
    render() {
        return (
          <div style={ { padding: 16 } }>
            <h4 style={ { marginTop: 32 } }>Classify Activity</h4>

            <RadioGroup onChange={ this.onChange } horizontal value={this.state.value}>
          <RadioButton value="none" pointColor="#999999">
            None
          </RadioButton>
          <RadioButton value="apple">
            Apple
          </RadioButton>
          <RadioButton value="orange">
            Orange
          </RadioButton>
        </RadioGroup>
          </div>
        );
      }
}

export default ChooseActivity;