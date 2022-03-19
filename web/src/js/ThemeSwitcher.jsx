import React, { Component } from 'react';

import ThemeLoader from './ThemeLoader';

import { ThemeContext } from './ThemeContext';

export default class ThemeSwitcher extends Component {
  render() {
    return (
      <ThemeContext.Consumer>
        {({ theme, validThemes, setTheme }) => (
          <div className="d-flex">
            <div className="dropdown mr-1">
              <button
                type="button"
                className="btn btn-secondary dropdown-toggle"
                id="dropdownMenuOffset"
                data-toggle="dropdown"
                aria-haspopup="true"
                aria-expanded="false"
                data-offset="10,20"
              >
                Theme
              </button>
              <div
                className="dropdown-menu"
                aria-labelledby="dropdownMenuOffset"
              >
                {validThemes.map((theme, index) => (
                  <a
                    key={index}
                    className="dropdown-item"
                    href="#"
                    onClick={() => setTheme(theme)}
                  >
                    {theme}
                  </a>
                ))}
              </div>
            </div>
          </div>
        )}
      </ThemeContext.Consumer>
    );
  }
}
