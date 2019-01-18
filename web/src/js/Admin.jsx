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
        <div>
          <a href="/api/auth/twitch/bot">Authenticate as bot</a> - <span>Use if you need to reauthenticate a bot below, or add a new one. You should probably copy-paste the link into incognito mode so you can log in as your bot account</span>
        </div>
        <table className="table table-sm">
          <thead>
            <tr>
              <th scope="col">Name</th>
              <th scope="col">Connected</th>
            </tr>
          </thead>
          <tbody>
          {this.state.bots.map((bot, index) =>
            <tr key={index}>
              <td>{bot.Name}</td>
              <td>{bot.Connected ? "Yes" : "No"}</td>
            </tr>
          )}
          </tbody>
        </table>
      </section>
    );
  }
}
