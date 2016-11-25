; (function(angular) {
	'use strict';

	angular.module('app.controllers').controller('PrivatesCtrl', ['$scope', 'Mikrotik', 'StrUtils',
		function($scope, Mikrotik, StrUtils) {
			$scope.clients = [];
			$scope.rowCollection = [];

			function requestClients() {
				Mikrotik.getClients().success(function(response) {
					var formatBytes = function(str) {
						if (!str) {
							return str;
						}

						var values = str.split('/');

						for (var i = 0; i < values.length; i++) {
							values[i] = StrUtils.formatBytes(values[i]);
						}

						return values.join('/');
					};

					$scope.clients = response.map(function(client) {
						client.target = client.target.split(',');
						client.maxLimit = formatBytes(client.maxLimit);
						client.burstLimit = formatBytes(client.burstLimit);
						client.burstThreshold = formatBytes(client.burstThreshold);

						$scope.rowCollection.push(client);

						return client;
					});

				});
			}

			$scope.reloadClients = function() {
				requestClients();
			};
		}
	]);
})(angular);