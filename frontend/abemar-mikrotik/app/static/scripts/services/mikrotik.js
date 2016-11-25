;(function(angular) {
	'use strict';

	angular.module('app.services').factory('Mikrotik', ['$http', 'Api',
		function($http, Api) {
			var service = {};

			service.getClients = function() {
				return $http({
					url: Api.getRoute('mikrotik/getClients/'),
					method: 'POST'
				});
			};

			return service;
		}
	]);
})(angular);
