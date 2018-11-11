import React, { Component } from "react";

import { ThemeContext } from "./ThemeContext";

export default class ThemeLoader extends Component {

  render() {
    return (
        <ThemeContext.Consumer>
          {({theme}) => (
            <link rel="stylesheet" href={this.themePath(theme)} />
            )}
        </ThemeContext.Consumer>
    );
  }

  themePath = (t) => {
    return '/static/themes/' + t + '/bootstrap.min.css';
  };
};
