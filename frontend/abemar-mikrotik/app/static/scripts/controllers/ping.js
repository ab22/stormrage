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
				var result = PingService.startPing($scope.ip);
				if (result === -1) {
					log('> Please enter an IP!');
					return;
				}

				log('> Request to ping [' + $scope.ip + '] sent...');
			};

			$scope.StopPing = PingService.stopPing;

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

				PingService.connect();
			}

			$scope.Reconnect = function() {
				if (!$scope.connected) {
					log('> Reconnecting...');
					connect();
				}
			};

			// Ctor
			$scope.$on('$viewContentLoaded', function() {
				log('> Connecting to server...');

				// If a ping process is running and the user switches views, then the service
				// will keep the previous instances of the controller's event functions and will
				// call those event functions when a new response is received from the server,
				// preventing all responses to be shown on the user's console.
				//
				// To avoid this, when the new controller is created, we setup the PingService
				// with the current functions of this instance.
				PingService.setup(onOpen, onClose, onMessage, onError);

				connect();
			});

			// Dtor
			$scope.$on('$destroy', function() {
				// In case user switches view without properly stopping an existing process,
				// we send a stop ping singal. If no process is running, the backend should do
				// nothing.
				$scope.StopPing();
			});
		}
	]);
})(angular);
