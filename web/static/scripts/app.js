class pb2ViewModel {
  constructor() {
    this.messages = ko.observableArray([]);
    this.reports = ko.observableArray([]);

    ws.subscribe('MessageReceived', (json) => {
      this.messages.push(json);
    });

    ws.subscribe('ReportReceived', (json) => {
      this.reports.push(json);
    });

    ws.subscribe('ReportHandled', (json) => {
      console.log('Report handled:', json);
      this.reports.remove((report) => {
        return report.ID == json.ReportID;
      });
    });
  }

  handleReport(action, report) {
    ws.publish('HandleReport', {
        'Action': action,
        'ChannelID': report.Channel.ID,
        'ReportID': report.ID,
      });
  }

  // Locally hides the report for this session only
  hideReport(report, el) {
    console.log(el.target.parent);
    console.log($(el.target).closest('.report').addClass('removing'));
    setTimeout(function() {
      vm.reports.remove(report);
    }, 500);
  }

  logIn() {
    console.log('log in');
    window.location.href = "/api/auth/twitch/user?redirect=/dashboard";
  }
}

class pb2WebSocket {
  constructor(nonce, user_id) {
    this.isOpen = false;
    this.socket = null;
    this.cbs = {};
    this.nonce = nonce;
    this.user_id = user_id;
  }

  subscribe(topic, cb) {
    if (typeof this.cbs[topic] !== 'undefined') {
      return;
    }

    this.cbs[topic] = cb;

    if (this.isOpen) {
      this.sendSubscribe(topic);
    }
  }

  publish(topic, data) {
    let payload = {
      'Type': 'Publish',
      'Topic': topic,
      'Data': data,
    };

    this.send(payload);
  }

  send(payload) {
    console.log('Sending', payload);
    this.socket.send(JSON.stringify(payload));
  }

  sendSubscribe(topic) {
    let payload = {
      'Type': 'Subscribe',
      'Topic': topic,
    };
    this.authorize(payload);
    this.send(payload);
  }

  authorize(payload) {
    if (this.nonce && this.user_id) {
      payload['Authorization'] = {
        'Nonce': this.nonce,
        'TwitchUserID': this.user_id,
      };
    }
  }

  connect() {
    if (this.isOpen) {
      // Already connected
      console.log('We are already connected????????');
      return;
    }

    console.log('Connecting to', ws_host);
    try {
      this.socket = new WebSocket(ws_host);
    } catch (e) {
      console.log("????????" + e);
      return;
    }
    console.log(this.socket);
    this.socket.binaryType = 'arraybuffer';
    this.socket.onopen = () => {
      console.log('Connected to PB2');
      this.isOpen = true;

      for (let topic in this.cbs) {
        this.sendSubscribe(topic);
      }
    };

    this.socket.onmessage = (e) => {
      if (typeof e.data !== 'string') {
        return;
      }

      let datas = e.data.split('\r\n');

      for (let i in datas) {
        let data = datas[i];
        var json = JSON.parse(data);
        if (typeof json !== 'object') {
          return;
        }

        if (typeof json['Type'] === 'string') {
          if (typeof json['Topic'] === 'string') {
            let cb = this.cbs[json['Topic']];
            if (cb !== undefined) {
              cb(json['Data']);
            }
          }
        }
      }
    };

    this.socket.onclose = (e) => {
      this.socket = null;
      this.isOpen = false;
      setTimeout(() => {
        this.connect();
      }, 500);
    }
  }
}

$(document).ready(function() {
  let nonce = window.localStorage.getItem('nonce');
  let user_id = window.localStorage.getItem('user_id');
  if (window.location.hash) {
    console.log('xd');
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

    history.replaceState('', document.title, window.location.pathname + window.location.search);
  }

  console.log('nonce:', nonce);

  // Read nonce
  console.log('Document ready');
  ws = new pb2WebSocket(nonce, user_id);
  vm = new pb2ViewModel();

  ws.connect();
  ko.applyBindings(vm);
});
