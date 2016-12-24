;(function(angular) {
	'use strict';

	angular.module('app.controllers').controller('LoginCtrl', ['$scope', '$location', 'ngToast', 'Auth',
		function($scope, $location, ngToast, Auth) {
			$scope.credentials = {
				identifier: '',
				password: ''
			};

			$scope.authenticate = function() {
				Auth.login($scope.credentials).then(
					function() {
						$location.path('/home');
					},

					function(response) {
						$scope.credentials.password = '';
						var message = response.data;

						ngToast.create({
							className: 'danger',
							content: message,
							dismissButton: true
						});
					}
				);
			};

		}
	]);
})(angular);
