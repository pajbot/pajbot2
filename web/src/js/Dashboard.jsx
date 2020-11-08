import React, { Component } from "react";
import WebSocketHandler from "./WebSocketHandler";
import { isLoggedIn } from "./auth";

const ReportActionUnknown = 0;
const ReportActionBan = 1;
const ReportActionTimeout = 2;
const ReportActionDismiss = 3;
const ReportActionUndo = 4;

const actionList = { ban: 'banned', unban: 'unbanned', timeout: 'timed out' };

function jsonifyResponse(response) {
  if (!response.ok) {
    throw response;
  }

  return response.json();
}

function parseError(response, onErrorParsed) {
  response
    .json()
    .then((obj) => onErrorParsed(obj))
    .catch((e) => console.error("Error parsing json from response:", e));
}

export default class Dashboard extends Component {
  constructor(props) {
    super(props);

    this.ws = new WebSocketHandler(props.element.getAttribute("data-wshost"));

    this.ws.subscribe("ReportHandled", (json) => {
      this.removeVisibleReport(json.ReportID);
    });

    let channel = null;

    try {
      channel = JSON.parse(props.element.getAttribute("data-extra"));
    } catch (e) {
      console.error("Error parsing channel JSON:", e);
    }

    if (channel) {
      if (channel.Bots == null) {
        channel.Bots = [];
      }
    }

    this.state = {
      bots: channel?.Bots || [],
      channel: channel?.Channel || "",
      currentBot: channel?.Bots[0] || null,
      channels: channel?.Channels || [],
      reports: [],
      userLookupLoading: false,
      userLookupData: null,
    };

    if (this.state.channel.ID) {
      console.log("Subscribe to report received", this.state.channel.ID);
      this.ws.subscribe(
        "ReportReceived",
        (json) => {
          this.addToVisibleReports(json);
        },
        {
          ChannelID: this.state.channel.ID,
        }
      );

      this.ws.connect();
    }
  }

  render() {
    if (!this.state.channel) {
      return (
        <section>
          <h2>Dashboards</h2>
          <p hidden={isLoggedIn()}>
            If you log in, you will get a list of channels you have access to
            here
          </p>
          <h4 hidden={!isLoggedIn()}>List of channel dashboards:</h4>
          <ul>
            {this.state.channels.map((channel, index) => (
              <li key={index}>
                <a href={`/c/${encodeURIComponent(channel.Name)}/dashboard`}>
                  {channel.Name} ({channel.ID})
                </a>
              </li>
            ))}
          </ul>
        </section>
      );
    }

    if (this.state.channel.ID == "") {
      return (
        <section>
          <h2>No bot is running in this channel</h2>
        </section>
      );
    }

    return (
      <section>
        <div
          className="alert alert-danger fade show"
          role="alert"
          hidden={!this.state.errorMessage}
        >
          {this.state.errorMessage}
          <button
            type="button"
            className="close"
            data-dismiss="alert"
            aria-label="Close"
          >
            <span aria-hidden="true">&times;</span>
          </button>
        </div>
        <table className="table table-sm">
          <thead>
            <tr>
              <th scope="col">Name</th>
              <th scope="col">Connected</th>
            </tr>
          </thead>
          <tbody>
            {this.state.bots.map((bot, index) => (
              <tr key={index}>
                <td>{bot.Name}</td>
                <td>{bot.Connected ? "Yes" : "No"}</td>
              </tr>
            ))}
          </tbody>
        </table>
        <div className="row dashboard">
          <div className="col reports">
            <h4>Reports</h4>
            {this.state.reports.map((report, index) => (
              <div className="report card mb-3" key={index}>
                <div className="card-header">
                  <i
                    className={
                      report.Channel.Type === "twitch"
                        ? "fab fa-twitch"
                        : "fab fa-discord"
                    }
                  />
                  &nbsp;{report.Channel.Name}&nbsp;
                  <span className="reporter">{report.Reporter.Name}</span>
                  &nbsp;reported&nbsp;
                  <span className="target">{report.Target.Name}</span>
                </div>
                <div className="card-body">
                  {report.Reason ? (
                    <span className="reason">{report.Reason}</span>
                  ) : null}
                  <a
                    target="_blank"
                    href={`https://logs.ivr.fi/?channel/${report.Channel.Name}/user/${report.Target.Name}`}
                  >
                    &nbsp;logs
                  </a>
                  <div>{report.Time}</div>
                  {report.Logs &&
                    report.Logs.map((value, key) => (
                      <div key={key}>{value}</div>
                    ))}
                  <button
                    className="card-link btn btn-danger"
                    title="Bans the user"
                    onClick={() =>
                      this.handleReport(
                        report.Channel.ID,
                        report.ID,
                        ReportActionBan
                      )
                    }
                  >
                    Ban
                  </button>
                  <button
                    className="card-link btn btn-danger"
                    title="Timeout the user"
                    onClick={() =>
                      this.handleReport(
                        report.Channel.ID,
                        report.ID,
                        ReportActionTimeout,
                        86400
                      )
                    }
                  >
                    Timeout 1d
                  </button>
                  <button
                    className="card-link btn btn-danger"
                    title="Timeout the user"
                    onClick={() =>
                      this.handleReport(
                        report.Channel.ID,
                        report.ID,
                        ReportActionTimeout,
                        604800
                      )
                    }
                  >
                    Timeout 7d
                  </button>
                  <button
                    className="card-link btn btn-danger"
                    title="Timeout the user"
                    onClick={() =>
                      this.handleReport(
                        report.Channel.ID,
                        report.ID,
                        ReportActionTimeout,
                        1209600
                      )
                    }
                  >
                    Timeout 14d
                  </button>
                  <button
                    className="card-link btn btn-danger"
                    title="Do nothing, but also don't untimeout the user"
                    onClick={() =>
                      this.handleReport(
                        report.Channel.ID,
                        report.ID,
                        ReportActionDismiss
                      )
                    }
                  >
                    Dismiss
                  </button>
                  <button
                    className="card-link btn btn-danger"
                    title="Undos timeout/ban"
                    onClick={() =>
                      this.handleReport(
                        report.Channel.ID,
                        report.ID,
                        ReportActionUndo
                      )
                    }
                  >
                    Undo
                  </button>
                  <button
                    className="card-link btn btn-light"
                    title="Hide this from your session"
                    onClick={() => this.removeVisibleReport(report.ID)}
                  >
                    Hide
                  </button>
                </div>
              </div>
            ))}
          </div>

          <div className="col userLookup">
            <h4>User lookup</h4>
            <form className="inline-group" onSubmit={this.lookupUser}>
              <div className="input-group input-group-sm mb-3">
                <input
                  type="text"
                  className="form-control"
                  placeholder="Recipient's username"
                  aria-label="Recipient's username"
                  aria-describedby="button-addon2"
                />
                <div className="input-group-append">
                  <input
                    type="submit"
                    className="btn btn-primary"
                    type="submit"
                    id="button-addon2"
                    value="Look up user"
                  />
                </div>
              </div>
            </form>
            <span hidden={this.state.userLookupLoading === false}>
              Loading...
            </span>
            {this.state.userLookupData &&
              (this.state.userLookupData.Actions.length > 0 ? (
                <div className="userData">
                  <span>
                    Listing latest{" "}
                    <strong>{this.state.userLookupData.Actions.length}</strong>{" "}
                    moderation actions on{" "}
                    <strong>{this.state.userLookupName}</strong>
                  </span>
                  <ul className="list-group">
                    {this.state.userLookupData.Actions.map((action, index) => (
                      <li className="list-group-item" key={index}>
                        <span>
                          [{action.Timestamp}] {action.UserName} {actionList[action.Action]}{" "}
                          {this.state.userLookupName} {action.Action === 'timeout' ? `for ${action.Duration}s: ` : ': '}
                          {action.Reason}
                        </span>
                      </li>
                    ))}
                  </ul>
                </div>
              ) : (
                "lol no bans"
              ))}
          </div>
        </div>
      </section>
    );
  }

  handleReport = (channelId, reportId, action, duration = null) => {
    let payload = {
      Action: action,
      ChannelID: channelId,
      ReportID: reportId,
    };

    if (action === "dismiss") {
      this.removeVisibleReport(reportId);
    }

    if (duration) {
      payload["Duration"] = duration;
    }

    this.ws.publish("HandleReport", payload);
  };

  removeVisibleReport = (id) => {
    const newReports = [];

    this.state.reports.map((report) => {
      if (report.ID !== id) {
        newReports.push(report);
      }
    });

    this.setState({
      ...this.state,
      reports: newReports,
    });
  };

  addToVisibleReports = (report) => {
    const newReports = this.state.reports;
    newReports.push(report);

    this.setState({
      ...this.state,
      reports: newReports,
    });
  };

  hasUserData = () => {
    return this.state.userLookupData !== null;
  };

  lookupUser = (e) => {
    e.preventDefault();

    let username = e.target[0].value;

    if (username.length < 2) {
      return;
    }

    this.setState({
      userLookupLoading: true,
      userLookupName: username,
    });

    fetch(
      "/api/channel/" +
        this.state.channel.ID +
        "/moderation/user?user_name=" +
        username
    )
      .then(jsonifyResponse)
      .then((myJson) => {
        this.setState({
          userLookupLoading: false,
          userLookupData: myJson,
        });
      })
      .catch((error) => {
        parseError(error, (e) => {
          this.setState({
            userLookupLoading: false,
            errorMessage: e.error,
          });
        });
      });
  };
}
