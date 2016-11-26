;(function(angular) {
	'use strict';

	angular.module('app.services').factory('PingService', ['Api',
		function(Api) {
			var pingService = {
				ws: null,
				onOpen: null,
				onClose: null,
				onMessage: null,
				onError: null
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

			pingService.setup = function(onOpen, onClose, onMessage, onError) {
				pingService.onOpen = onOpen;
				pingService.onClose = onClose;
				pingService.onMessage = onMessage;
				pingService.onError = onError;
			};

			pingService.connect = function() {
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
					if (pingService.onOpen) {
						pingService.onOpen(evt);
					}
				};

				ws.onclose = function(evt) {
					pingService.ws = null;

					if (pingService.onClose) {
						pingService.onClose(evt);
					}
				};

				ws.onmessage = function(evt) {
					if (pingService.onMessage) {
						pingService.onMessage(evt);
					}
				};

				ws.onerror = function(evt) {
					if (pingService.onError) {
						pingService.onError(evt);
					}
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

			pingService.startPing = function(ip) {
				ip = ip.trim();

				if (!ip) {
					return -1;
				}

				var request = {
					option: 0,
					ip: ip
				};

				pingService.sendJSON(request);
				return 0;
			};

			pingService.stopPing = function() {
				var request = {
					option: 1
				};

				pingService.sendJSON(request);
			};

			return pingService;
		}
	]);
})(angular);
