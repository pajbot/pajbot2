import React, { Component } from "react";
import WebSocketHandler from "./WebSocketHandler"

export default class Dashboard extends Component {

	constructor(props) {
		super(props);

		const auth = this.loadAuth();

		this.state = {
			nonce: auth.nonce,
			userId: auth.userId,
			reports: [],
		};

		this.ws = new WebSocketHandler(props.wshost, auth.nonce, auth.userId);

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
				<button className="btn btn-twitch" hidden={this.isLoggedIn()} onClick={this.logIn}><i className="fab fa-twitch" /> Connect with Twitch</button>
				<button className="btn btn-twitch" hidden={!this.isLoggedIn()} onClick={this.logOut}>Log out</button>
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

	isLoggedIn = () => {
		return this.state.nonce !== null && this.state.userId !== null;
	}

	loadAuth = () => {
		let nonce = window.localStorage.getItem('nonce');
		let user_id = window.localStorage.getItem('user_id');
		if (window.location.hash) {
			let hash = window.location.hash.substring(1);
			let hashParts = hash.split(';');
			let hashMap = {};
			for (let i = 0; i < hashParts.length; ++i) {
				let p = hashParts[i].split('=');
				if (p.length == 2) {
					hashMap[p[0]] = p[1];
				}
			}

			if ('nonce' in hashMap && 'user_id' in hashMap) {
				nonce = hashMap['nonce'];
				user_id = hashMap['user_id'];
				window.localStorage.setItem('nonce', nonce);
				window.localStorage.setItem('user_id', user_id);
			}

			history.replaceState(
				'', document.title,
				window.location.pathname + window.location.search);
		}

		return { nonce: nonce, userId: user_id };
	}

	logIn = () => {
		window.location.href = '/api/auth/twitch/user?redirect=/dashboard';
	}

	logOut = () => {
		this.setState({
			...this.state,
			nonce: null,
			userId: null,
			reports: [],
		});

		window.localStorage.removeItem('nonce');
		window.localStorage.removeItem('user_id');
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