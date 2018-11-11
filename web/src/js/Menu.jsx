import React, { Component } from "react";
import LogInButton from "./LogInButton"

export default class Menu extends Component {

  constructor(props) {
    super(props);

    const menuItems = [
      {
        link: "/",
        name: "Home",
      },
      {
        link: "/dashboard",
        name: "Dashboard",
      },
    ];

    this.state = {
      menuItems: menuItems,
    };
  }

  render() {
    return (
      <nav className="navbar navbar-expand-lg navbar-light bg-light">
        <a className="navbar-brand" href="#">pajbot2</a>
        <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNavAltMarkup" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
          <span className="navbar-toggler-icon"></span>
        </button>
        <div className="collapse navbar-collapse" id="navbarNavAltMarkup">
          <div className="navbar-nav">
{this.state.menuItems.map((menuItem, index) =>
            <a key={index} className={"nav-item nav-link "+(window.location.pathname == menuItem.link ? "active" : "")} href={menuItem.link}>{menuItem.name}</a>
          )}
          </div>
        </div>
        <LogInButton />
      </nav>
    );
  }

}
