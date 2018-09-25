class pb2ViewModel {
  constructor() {
    this.reports = ko.observableArray([]);
    this.nonce = ko.observable(null);
    this.user_id = ko.observable(null);

    this.loggedIn = ko.computed(() => {
      return this.nonce() != null && this.user_id() != null;
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

    this.loadAuth();
  }

  loadAuth() {
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

    this.nonce(nonce);
    this.user_id(user_id);
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
    window.location.href = '/api/auth/twitch/user?redirect=/dashboard';
  }

  logOut() {
    this.nonce(null);
    this.user_id(null);
    window.localStorage.removeItem('nonce');
    window.localStorage.removeItem('user_id');
  }

  onConnected() {
    this.reports.removeAll();
  }

  makeLogLink(channelName, targetName) {
    return 'https://api.gempir.com/channel/' + channelName + '/user/' + targetName;
  }
}

class pb2WebSocket {
  constructor() {
    this.isOpen = false;
    this.socket = null;
    this.cbs = {};
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
    if (vm.nonce() && vm.user_id()) {
      payload['Authorization'] = {
        'Nonce': vm.nonce(),
        'TwitchUserID': vm.user_id(),
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
      console.log('????????' + e);
      return;
    }
    console.log(this.socket);
    this.socket.binaryType = 'arraybuffer';
    this.socket.onopen = () => {
      vm.onConnected();
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
  ws = new pb2WebSocket();
  vm = new pb2ViewModel();

  ws.connect();
  ko.applyBindings(vm);
});
