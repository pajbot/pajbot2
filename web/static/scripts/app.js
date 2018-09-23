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

  send(msg) {
    console.log('Sending', msg);
    this.socket.send(JSON.stringify(msg));
  }

  sendSubscribe(topic) {
    this.send({'Type': 'Subscribe', 'Topic': topic});
  }

  connect() {
    if (this.isOpen) {
      // Already connected
      console.log('We are already connected????????');
      return;
    }

    this.socket = new WebSocket(ws_host);
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
  console.log('Document ready');
  ws = new pb2WebSocket();
  vm = new pb2ViewModel();

  console.log('Call connect on ws');
  ws.connect();
  ko.applyBindings(vm);
});
