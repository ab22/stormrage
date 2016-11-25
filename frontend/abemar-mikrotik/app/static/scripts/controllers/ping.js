; (function(angular) {
	'use strict';

	angular.module('app.controllers').controller('PingCtrl', ['$scope',
		function($scope) {
			$scope.ws = null;
			var logObj = $('.log');
			$scope.ip = '';

			function appendToLog(msg) {
				var d = logObj[0];
				var doScroll = d.scrollTop === d.scrollHeight - d.clientHeight;
				msg.appendTo(logObj);
				if (doScroll) {
					d.scrollTop = d.scrollHeight - d.clientHeight;
				}
			}

			function log(msg) {
				appendToLog($('<div/>').text(msg));
			}

			function clearLog() {
				logObj.empty();
			}

			function onConnect() {
				log('> Connection Established!');
			}

			function onClose() {
				$scope.ws = null;
				$scope.$apply();
				log('> Connection Closed!');
			}

			function onError() {
				log('> An error ocurred!');
			}

			function onMessage(evt) {
				var data = evt.data.replace(/[\u0000-\u0019]+/g, '');
				var response = JSON.parse(data);

				if (response.error) {
					log(response.error);
					return;
				}

				if (response.payload) {
					log(response.payload);
					return;
				}

				log(evt.data);
			}

			$scope.StartPing = function() {
				var ip = $scope.ip.trim();
				if (!$scope.ip) {
					log('> Please enter an IP!');
					return;
				}

				var request = {
					option: 0,
					ip: ip
				};

				$scope.ws.send(JSON.stringify(request));
				log('> Request to ping [' + ip + '] sent...');
			};

			$scope.StopPing = function() {
				var request = { option: 1 };
				$scope.ws.send(JSON.stringify(request));
			};

			$scope.ClearLog = function() {
				clearLog();
			};

			function createWebSocket() {
				var host = window.location.host;
				var url = 'ws://' + host + '/api/ws/onConnect/';

				$scope.ws = new WebSocket(url);
				$scope.ws.onopen = onConnect;
				$scope.ws.onclose = onClose;
				$scope.ws.onmessage = onMessage;
				$scope.ws.onerror = onError;
			}

			$scope.Reconnect = function() {
				if (!$scope.ws) {
					log('> Reconnecting...');
					createWebSocket();
				}
			};

			function init() {
				log('> Connecting to server...');
				createWebSocket();
			}

			init();
		}
	]);
})(angular);
