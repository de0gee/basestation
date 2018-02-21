import * as React from 'react';
import { RadioGroup, RadioButton } from 'react-radio-buttons'; 

class ChooseActivity extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            radioprops: [
                {label: 'param1', value: 0 },
                {label: 'param2', value: 1 }
              ],
        };
      }

      onChange(value) {
        console.log(value);
      }
    
    render() {
        return (
          <div>
              <strong>Classify activity:</strong>
            <RadioGroup onChange={ this.onChange } horizontal>
          <RadioButton value="none">
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