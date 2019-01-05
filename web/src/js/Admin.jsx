import React, { Component } from "react";
import WebSocketHandler from "./WebSocketHandler";

export default class Admin extends Component {

  constructor(props) {
    super(props);

    let extra = JSON.parse(props.element.getAttribute('data-extra'));

    this.state = {
      bots: extra.Bots,
    };
  }

  render() {
    return (
      <section>
        <h4>Admin</h4>
        <table className="table table-sm">
          <thead>
            <tr>
              <th scope="col">Name</th>
            </tr>
          </thead>
          <tbody>
          {this.state.bots.map((bot, index) =>
            <tr key={index}>
              <td>{bot.Name}</td>
            </tr>
          )}
          </tbody>
        </table>
      </section>
    );
  }
}
