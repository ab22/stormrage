; (function(angular) {
	'use strict';

	angular.module('app.controllers').controller('PingCtrl', ['$scope', 'Auth', 'PingService',
		function($scope, Auth, PingService) {
			$scope.connected = false;
			$scope.ip = '';
			var logObj = $('.log');

			function appendToLog(msg) {
				var d = logObj[0];
				var doScroll = d.scrollTop === d.scrollHeight - d.clientHeight;
				msg.appendTo(logObj);
				if (doScroll) {
					d.scrollTop = d.scrollHeight - d.clientHeight;
				}
			}

			function log(msg) {
				appendToLog($('<div/>').html(msg));
			}

			function clearLog() {
				logObj.empty();
			}

			function onOpen() {
				$scope.connected = true;
				$scope.$apply();
				log('> Connection Established!');
			}

			function onClose() {
				$scope.connected = false;
				$scope.$apply();

				log('> Connection Closed!');
			}

			function onError() {
				Auth.checkAuthentication();
				log('> An error ocurred!');
			}

			function onMessage(evt) {
				// Remove all new line characters (\n) or any other non supported JSON char to
				// prevent JSON.parse from throwing an error.
				var data = evt.data
								.replace(/\n/g, '<br />')
								.replace(/[\u0000-\u0019]+/g, '');
				var response;

				try {
					response = JSON.parse(data);
				} catch (err) {
					log('> Error parsing data from server!');
					log('> ' + err);
					PingService.disconnect();
					return;
				}

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

				PingService.sendJSON(request);
				log('> Request to ping [' + ip + '] sent...');
			};

			$scope.StopPing = function() {
				var request = { option: 1 };
				PingService.sendJSON(request);
			};

			$scope.ClearLog = function() {
				clearLog();
			};

			function connect() {
				var state = PingService.getState();

				if (state === -1) {
					log('> WebSockets are not supported by this browser!');
					return;
				}

				if (state <= 1) {
					$scope.connected = true;
					log('> Connection already established!');
					return;
				}

				PingService.connect(onOpen, onClose, onMessage, onError);
			}

			$scope.Reconnect = function() {
				if (!$scope.connected) {
					log('> Reconnecting...');
					connect();
				}
			};

			function init() {
				log('> Connecting to server...');
				connect();
			}

			init();
		}
	]);
})(angular);
