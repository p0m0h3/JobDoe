const WebSocket = require('ws');

const address = 'ws://127.0.0.1:7001/v1/task/4d0ac7003ad43f51c43794ab88b6c09072389b875a89adfefafef16dd7bfbf83/stream'
const websocket = new WebSocket(address, {
	perMessageDeflate: true,
});

websocket.on("message", function incoming(data) {
	console.log(data.toString())
});
