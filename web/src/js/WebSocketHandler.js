export default class WebSocketHandler {
  constructor(wsHost) {
    this.isOpen = false;
    this.wsHost = wsHost;
    this.socket = null;
    this.cbs = {};
  }

  subscribe(topic, cb, params) {
    console.log('Params:', params);
    if (typeof this.cbs[topic] !== 'undefined') {
      return;
    }

    this.cbs[topic] = {
      cb: cb,
      params: params,
    };

    if (this.isOpen) {
      this.sendSubscribe(topic, params);
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
    this.socket.send(JSON.stringify(payload));
  }

  sendSubscribe(topic, params) {
    console.log('Params2:', params);
    let payload = {
      'Type': 'Subscribe',
      'Topic': topic,
      'Data': params,
    };

    console.log('Sending', payload);

    this.send(payload);
  }

  connect() {
    if (this.isOpen) {
      // Already connected
      console.log('We are already connected????????');
      return;
    }

    console.log('Connecting to', this.wsHost);
    try {
      this.socket = new WebSocket(this.wsHost);
    } catch (e) {
      console.log('????????' + e);
      return;
    }
    console.log(this.socket);
    this.socket.binaryType = 'arraybuffer';
    this.socket.onopen = () => {
      console.log('Connected to PB2');
      this.isOpen = true;

      for (let topic in this.cbs) {
        this.sendSubscribe(topic, this.cbs[topic].params);
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
            let cb = this.cbs[json['Topic']].cb;
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
