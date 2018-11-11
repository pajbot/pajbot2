import React, { Component } from "react";

import {isLoggedIn, logIn, logOut} from "./auth"

export default class LogInButton extends Component {

  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div>
      <button className="btn btn-twitch" hidden={isLoggedIn()} onClick={logIn}><i className="fab fa-twitch" /> Connect with Twitch</button>
      <button className="btn btn-twitch" hidden={!isLoggedIn()} onClick={logOut}>Log out</button>
    </div>
    );
  }
}
