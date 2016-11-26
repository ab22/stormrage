;(function(angular) {
	'use strict';

	angular.module('app.services').factory('PingService', ['Api',
		function(Api) {
			var pingService = {
				ws: null
			};

			// getState returns the current websocket.readyState.
			//
			// STATES:
			//         0 - CONNECTING
			//         1 - CONNECTED
			//         2 - CLOSING
			//         3 - CLOSED
			pingService.getState = function() {
				if (!window.WebSocket) {
					return -1;
				}

				var ws = pingService.ws;

				// No websocket instance. Assume it's closed.
				if (!ws) {
					return 3;
				}

				// No readyState attribute found on instance. Assume it's closed.
				if (!ws.readyState && ws.readyState !== 0) {
					return 3;
				}

				return ws.readyState;
			};

			pingService.connect = function(onOpen, onClose, onMessage, onError) {
				var state = pingService.getState();

				// There's alread an active websocket connection open or websockets are not
				// supported.
				if (state <= 1) {
					return state;
				}

				var host = window.location.host + '/';
				var url = 'wss://' + host + Api.getRoute('ws/onConnect/');
				var ws = new WebSocket(url);

				ws.onopen = function(evt) {
					onOpen(evt);
				};

				ws.onclose = function(evt) {
					pingService.ws = null;
					onClose(evt);
				};

				ws.onmessage = function(evt) {
					onMessage(evt);
				};

				ws.onerror = function(evt) {
					onError(evt);
				};

				pingService.ws = ws;
				return 0;
			};

			pingService.disconnect = function() {
				if (!pingService.ws) {
					return;
				}

				pingService.ws.close();
				pingService.ws = null;
			};

			pingService.sendJSON = function(request) {
				if (!pingService.ws) {
					return;
				}

				var query = JSON.stringify(request);
				pingService.ws.send(query);
			};

			pingService.sendText = function(text) {
				if (!pingService.ws) {
					return;
				}

				pingService.ws.send(text);
			};

			return pingService;
		}
	]);
})(angular);
