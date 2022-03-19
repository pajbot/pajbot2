import React, { Component } from 'react';
import WebSocketHandler from './WebSocketHandler';

export default class Commands extends Component {
  constructor(props) {
    super(props);

    let extra = JSON.parse(props.element.getAttribute('data-extra'));

    this.state = {
      commands: extra,
    };

    console.log(this.state);
  }

  render() {
    return (
      <section>
        <h4>Commands</h4>
        <table className="table table-sm">
          <thead>
            <tr>
              <th scope="col">Name</th>
              <th scope="col">Description</th>
            </tr>
          </thead>
          <tbody>
            {this.state.commands.map((command, index) => (
              <tr key={index}>
                <td>{command.Name}</td>
                <td>{command.Description}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>
    );
  }
}
