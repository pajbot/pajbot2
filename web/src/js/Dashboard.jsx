import React, { Component } from "react";
import WebSocketHandler from "./WebSocketHandler";
import {getAuth} from "./auth"

export default class Dashboard extends Component {

	constructor(props) {
		super(props);

		this.state = {
			reports: [],
		};

		this.ws = new WebSocketHandler(props.wshost, getAuth());

		this.ws.subscribe('ReportReceived', (json) => {
			console.log("Report Inc:", json);
			this.addToVisibleReports(json);
		});

		this.ws.subscribe('ReportHandled', (json) => {
			console.log('Report handled:', json);
			this.removeVisibleReport(json.ReportID);
		});

		this.ws.connect();
	}

	render() {
		return (
			<div className="dashboard">
				<div className="reports">
					{this.state.reports.map((report, index) =>
						<div className="report card mb-3" key={index}>
							<div className="card-header">
								<i className={report.Channel.Type === "twitch" ? "fab fa-twitch" : "fab fa-discord"} />&nbsp;{report.Channel.Name}&nbsp;
								<span className="reporter">{report.Reporter.Name}</span>&nbsp;reported&nbsp;<span className="target">{report.Target.Name}</span>
							</div>
							<div className="card-body">
								{report.Reason ? <span className="reason">{report.Reason}</span> : null}
								<a target="_blank" href={`https://api.gempir.com/channel/${report.Channel.Name}/user/${report.Target.Name}`}>&nbsp;logs</a>
								<div>{report.Time}</div>
								{report.Logs ? report.Logs.map((value, key) => 
									<div key={key}>{value}</div>	
								) : null}
								<button className="card-link btn btn-danger" title="Bans the user" onClick={() => this.handleReport(report.Channel.ID, report.ID, "ban")}>Ban</button>
								<button className="card-link btn btn-danger" title="Bans the user" onClick={() => this.handleReport(report.Channel.ID, report.ID, "timeout", 86400)}>Timeout 1d</button>
								<button className="card-link btn btn-danger" title="Bans the user" onClick={() => this.handleReport(report.Channel.ID, report.ID, "timeout", 604800)}>Timeout 7d</button>
								<button className="card-link btn btn-danger" title="Bans the user" onClick={() => this.handleReport(report.Channel.ID, report.ID, "timeout", 1209600)}>Timeout 14d</button>
								<button className="card-link btn btn-danger" title="Bans the user" onClick={() => this.handleReport(report.Channel.ID, report.ID, "dismiss")}>Dismiss</button>
								<button className="card-link btn btn-danger" title="Undos timeout/ban" onClick={() => this.handleReport(report.Channel.ID, report.ID, "undo")}>Undo</button>
								<button className="card-link btn btn-light" title="Hide this from your session" onClick={() => this.removeVisibleReport(report.ID)}>Hide</button>
							</div>
						</div>

					)}
				</div>
			</div>
		);
	}

	handleReport = (channelId, reportId, action, duration = null) => {
		let payload = {
			'Action': action,
			'ChannelID': channelId,
			'ReportID': reportId,
		};

		if (action === "dismiss") {
			this.removeVisibleReport(reportId);
		}

		if (duration) {
			payload['Duration'] = duration;
		}

		this.ws.publish('HandleReport', payload);
	}

	removeVisibleReport = (id) => {
		const newReports = [];
		
		this.state.reports.map(report => {
			if (report.ID !== id){
				newReports.push(report);
			}  	
		})

		this.setState({
			...this.state,
			reports: newReports
		});
	}

	addToVisibleReports = (report) => {

		const newReports = this.state.reports;
		newReports.push(report);

		this.setState({
			...this.state,
			reports: newReports,
		})
	}
}
