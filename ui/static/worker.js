function log(message) {
	console.debug(
		`[worker.js @ ${new Date().toISOString().slice(0, 19).replace('T', ' ')}] ${message}`
	);
}


const Status = {
	UNINITIALIZED: 'uninitialized',
	CONNECTING: 'connecting',
	OPEN: 'open',
	ERRORED: 'errored',
	CLOSED: 'closed',
	TIMEOUT: 'timeout'
};
let wsStatus = Status.UNINITIALIZED;
let socket;

let statusListeners = [];
function broadcastStatus(newStatus) {
	log(`WebSocket status changed: [${newStatus}]`);
	statusListeners.forEach((listener) => listener(newStatus));
	wsStatus = newStatus;
}

let message_listeners = [];
function broadcastMessage(event) {
	log(`Received message: ${event.data}`.slice(0, 250));
	message_listeners.forEach((listener) => listener(event.data));
}

let messageQueue = [];
function send(data) {
	if (wsStatus !== Status.OPEN) {
		log(`[ws: ${wsStatus}] Queuing message: [${data}]`);
		messageQueue.push(data);
		return;
	}
	log(`[ws: ${wsStatus}] Sending message: [${data}]`);
	socket.send(data);
}

function clearQueue() {
	console.log(
		`Clearing message queue: [${messageQueue.length} entr${
			messageQueue.length === 1 ? 'y' : 'ies'
		}}]`
	);
	messageQueue.forEach((data) => send(data));
	messageQueue = [];
}

let connectionAttempts = 0;
function connect() {
	connectionAttempts++;
	log(`Attempting to connect to WebSocket... [Attempts: ${connectionAttempts}]`);
	broadcastStatus(Status.CONNECTING);
	socket = new WebSocket(`ws://${self.location.hostname}/api/instance/default/socket`);

	let connectionTimeout = setTimeout(() => {
		log({ msg: "Closing socket if it still hasn't opened", readyState: socket.readyState });
		if (socket.readyState !== WebSocket.OPEN) {
			socket.close();
		}
	}, 5000);

	socket.addEventListener('open', () => {
		clearTimeout(connectionTimeout);
		broadcastStatus(Status.OPEN);
		clearQueue();
	});

	socket.addEventListener('close', () => {
		broadcastStatus(Status.CLOSED);
		retry();
	});
	socket.addEventListener('error', () => broadcastStatus(Status.ERRORED));
	socket.addEventListener('message', broadcastMessage);
}

function retry() {
	broadcastStatus(Status.TIMEOUT);
	setTimeout(connect, 1000);
}

connect();


self.onconnect = function (e) {
	const port = e.ports[0];

	port.postMessage(['status', wsStatus]);

	statusListeners.push((newStatus) => {
		port.postMessage(['status', newStatus]);
	});

	message_listeners.push((data) => {
		port.postMessage(['message', data]);
	});

	port.onmessage = function (e) {
		const [action, data] = e.data;
		switch (action) {
			case 'send':
				send(data);
		}
	};

	port.start();
};
