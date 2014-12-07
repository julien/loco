(function () {

  var ws = new WebSocket('ws://localhost:8000/ws');
  ws.onmessage = function (e) {
    document.location.reload();
  };
}())
