import React, { Component } from 'react';

import { isLoggedIn, logIn, logOut, loggedInUsername } from './auth';
import { readCookie } from './cookie';

export default class LogInButton extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div>
        <span hidden={!isLoggedIn()}>
          Logged in as <a href="/profile">{loggedInUsername()}</a>
        </span>
        &emsp;
        <button
          className="btn btn-twitch"
          hidden={isLoggedIn()}
          onClick={logIn}
        >
          <i className="fab fa-twitch" /> Connect with Twitch
        </button>
        <button
          className="btn btn-twitch"
          hidden={!isLoggedIn()}
          onClick={logOut}
        >
          Log out
        </button>
      </div>
    );
  }
}
