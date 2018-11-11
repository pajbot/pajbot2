import React, { Component } from "react";

import { ThemeContext } from "./ThemeContext";

import { createCookie, readCookie } from './cookie';

export default class ThemeProvider extends Component {


  validThemes = ['default', 'dark'];

  constructor(props) {
    super(props);

    let savedTheme = readCookie('currentTheme');

    if (!savedTheme || this.validThemes.indexOf(savedTheme) === -1) {
      savedTheme = 'default';
    }

    this.state = {
      theme: savedTheme,
    };
  }

  render() {
    return (
      <ThemeContext.Provider value={{
        theme: this.state.theme,
        validThemes: this.validThemes,
        setTheme: this.setTheme,
      }}>
      {this.props.children}
    </ThemeContext.Provider>
    );
  }

  setTheme = (t) => {
    this.setState({
      theme: t,
    });

    // Save theme
    createCookie('currentTheme', t, 3600);
  }
};
