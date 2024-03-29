import React, { Component } from 'react';
import LogInButton from './LogInButton';
import ThemeSwitcher from './ThemeSwitcher';
import { isLoggedIn } from './auth';

export default class Menu extends Component {
  constructor(props) {
    super(props);

    this.menuItems = [
      {
        link: '/',
        name: 'Home',
        requireLogin: false,
      },
      {
        link: '/admin',
        name: 'Admin',
        requireLogin: true,
      },
      {
        link: '/dashboard',
        name: 'Dashboard',
        requireLogin: true,
      },
    ];
  }

  render() {
    return (
      <nav className="navbar navbar-expand-lg navbar-dark bg-dark">
        <a className="navbar-brand" href="#">
          pajbot2
        </a>
        <button
          className="navbar-toggler"
          type="button"
          data-toggle="collapse"
          data-target="#navbarNavAltMarkup"
          aria-controls="navbarNavAltMarkup"
          aria-expanded="false"
          aria-label="Toggle navigation"
        >
          <span className="navbar-toggler-icon"></span>
        </button>
        <div className="collapse navbar-collapse" id="navbarNavAltMarkup">
          <div className="navbar-nav">
            {this.menuItems
              .filter(
                (item) =>
                  !item.requireLogin || (item.requireLogin && isLoggedIn())
              )
              .map((menuItem, index) => (
                <a
                  key={index}
                  className={`nav-item nav-link ${
                    window.location.pathname == menuItem.link ? 'active' : ''
                  }`}
                  href={menuItem.link}
                >
                  {menuItem.name}
                </a>
              ))}
          </div>
        </div>
        <ThemeSwitcher ThemeContext={this.props.themeContext} />
        <LogInButton />
      </nav>
    );
  }
}
