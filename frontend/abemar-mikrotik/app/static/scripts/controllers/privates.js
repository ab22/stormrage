; (function(angular) {
	'use strict';

	angular.module('app.controllers').controller('PrivatesCtrl', ['$scope', 'Mikrotik', 'StrUtils',
		function($scope, Mikrotik, StrUtils) {
			$scope.clients = [];
			$scope.rowCollection = [];

			function requestClients() {
				Mikrotik.getClients().then(
					function(response) {
						var data = response.data || [];
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

						$scope.clients = data.map(function(client) {
							client.target = client.target.split(',');
							client.maxLimit = formatBytes(client.maxLimit);
							client.burstLimit = formatBytes(client.burstLimit);
							client.burstThreshold = formatBytes(client.burstThreshold);

							return client;
						});

						$scope.rowCollection = $scope.clients;
					}
				);
			}

			$scope.reloadClients = function() {
				requestClients();
			};
		}
	]);
})(angular);
