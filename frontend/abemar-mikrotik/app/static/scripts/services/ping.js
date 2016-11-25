;(function(angular) {
	'use strict';

	angular.module('app.services').factory('PingService', ['Api',
		function(Api) {
			var pingService = {
				ws: null
			};

			pingService.connect = function(onOpen, onClose, onMessage, onError) {
				if (!window.WebSocket) {
					return -1;
				}

				if (pingService.ws) {
					return 0;
				}

				var host = window.location.host + '/';
				var url = 'ws://' + host + Api.getRoute('ws/onConnect/');
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
