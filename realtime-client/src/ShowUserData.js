import * as React from 'react';
import axios from 'axios';

class ShowUserData extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
          username: "",
       };
      var self = this;
      axios.get(window.location.href.replace(":3000",":8002")+'username')
        .then(function (response) {
          console.log(response);
          if (response.data.success === true) {
            self.setState({username: response.data.message});
          } else {
          }
        })
        .catch(function (error) {
          console.log(error);
        });
      }
    
    render() {
        return (
          <span>
            <h2>Username: <h1>{this.state.username}</h1></h2>
          </span>
        );
      }
}

export default ShowUserData;