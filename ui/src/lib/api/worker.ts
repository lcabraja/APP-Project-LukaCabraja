import { Connection } from '$lib/stores/ConnectionState';
import { LiveScoreState, translateState } from '$lib/stores/HandballState';

function base64ToUtf8(jsonString: string) {
	return new TextDecoder().decode(Uint8Array.from(atob(jsonString), (c) => c.charCodeAt(0)));
}

export function createWorker(): SharedWorker {
	const wsWorker = new SharedWorker('/worker.js');
	wsWorker.port.start();

	wsWorker.port.onmessage = function (e) {
		console.debug(`Message received from worker: [Length: ${e.data.toString().length}]`);
		switch (e.data[0]) {
			case 'status':
				switch (e.data[1]) {
					case 'connecting':
						Connection.set({
							status: 'connecting',
							initials: 'CO',
							fill: 'fill-black',
							background: 'bg-green-300'
						});
						break;
					case 'open':
						Connection.set({
							status: 'open',
							initials: 'OK',
							fill: 'fill-white',
							background: 'bg-green-600'
						});
						break;
					case 'closed':
						Connection.set({
							status: 'closed',
							initials: 'CL',
							fill: 'fill-white',
							background: 'bg-red-600'
						});
						break;
					case 'error':
						Connection.set({
							status: 'error',
							initials: 'ER',
							fill: 'fill-white',
							background: 'bg-orange-600'
						});
						break;
					case 'timeout':
						Connection.set({
							status: 'timeout',
							initials: 'TO',
							fill: 'fill-white',
							background: 'bg-purple-600'
						});
						break;
				}
				break;
			case 'message':
				try {
					let json = JSON.parse(e.data[1]);
					switch (json['action']) {
						case 'onboard':
							const innerJson = JSON.parse(base64ToUtf8(json.data));
							const liveScoreState = translateState(innerJson);
							LiveScoreState.set(liveScoreState);
							console.log({ liveScoreState });
							break;
					}
				} catch (err) {
					console.error(`Error recieving message: `, err, e.data[1]);
				}
				break;
		}
	};

	return wsWorker;
}
