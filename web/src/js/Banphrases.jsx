import React, { Component } from "react";
import WebSocketHandler from "./WebSocketHandler";

const ReportActionUnknown = 0;
const ReportActionBan     = 1;
const ReportActionTimeout = 2;
const ReportActionDismiss = 3;
const ReportActionUndo    = 4;

function jsonifyResponse(response) {
  if (!response.ok) {
    throw response;
  }

  return response.json();
}

function parseError(response, onErrorParsed) {
  response.json()
    .then(obj => onErrorParsed(obj))
    .catch(e => console.error('Error parsing json from response:', e));
}

export default class Banphrases extends Component {

  constructor(props) {
    super(props);

    console.log('props:', props);

    let banphrases = [
      {
        'id': 1,
        'description': 'lol description',
        'phrase': 'lol phrase',
        'enabled': true,
      },
      {
        'id': 2,
        'description': 'lol description',
        'phrase': 'lol phrase',
        'enabled': true,
      },
      {
        'id': 3,
        'description': 'lol description',
        'phrase': 'lol phrase',
        'enabled': false,
      },
      {
        'id': 4,
        'description': 'lol description',
        'phrase': 'lol phrase',
        'enabled': true,
      },
    ];

    this.state = {
      banphrases: banphrases,
    };

    // loadBanphrases();
  }

  render() {
    return (
      <section>
        <h4>Banphrases</h4>
        <table className="table table-sm">
          <thead>
            <tr>
              <th scope="col">ID</th>
              <th scope="col">Description</th>
              <th scope="col">Phrase</th>
              <th scope="col">Enabled</th>
              <th scope="col"></th>
            </tr>
          </thead>
          <tbody>
          {this.state.banphrases.map((bp, index) =>
            <tr key={index}>
              <td>{bp.id}</td>
              <td>{bp.description}</td>
              <td>{bp.phrase}</td>
              <td><input type="checkbox" checked={bp.enabled} /></td>
              <td>button</td>
            </tr>
          )}
          </tbody>
        </table>
      </section>
    );
  }
}
