window.onload = function () {

  var ws = new WebSocket('ws://localhost:8000/ws');

  ws.onclose = function (e) {
    console.log('on close');
  };

  ws.onmessage = function (e) {
    console.log('on message', e);
    document.location.reload();
  };


};
