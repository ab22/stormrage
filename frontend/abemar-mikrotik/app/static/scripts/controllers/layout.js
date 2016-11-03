;(function(angular) {
	'use strict';

	angular.module('app.controllers').controller('MainLayoutCtrl', ['$scope','$location', 'Auth',
		function($scope, $location, Auth) {
			$scope.window = {};
			$scope.isCollapsed = true;

			$scope.signOut = function() {
				Auth.logout().success(function() {
					$location.path('/login');
				});
			};

			$scope.activeOption = null;
			$scope.topMenu = [];
			$scope.menu = [
				[
					{
						label: 'Inicio',
						icon: 'fa-home',
						link: '/main/home',
						responsiveOnly: false,
					}
				],
				[
					{
						label: 'Cerrar Sesión',
						icon: 'fa-sign-out',
						link: '',
						onClick: $scope.signOut,
						responsiveOnly: false,
					}
				]
			];

			function hideResponsiveMenu() {
				if ($scope.isResponsiveMode()) {
					$scope.isCollapsed = !$scope.isCollapsed;
				}
			}

			function setActiveOption(option) {
				$scope.activeOption = option;
			}

			$scope.optionOnClick = function(option) {
				hideResponsiveMenu();

				if (typeof option.onClick !== 'undefined') {
					option.onClick();
					return;
				}

				setActiveOption(option);
			};

			$scope.isResponsiveMode = function() {
				return $scope.window.width <= 767;
			};

			$scope.showResponsiveOption = function(option) {
				var showResponsive = !option.responsiveOnly || (option.responsiveOnly && $scope.isResponsiveMode());

				return showResponsive;
			};

			function determineActiveOption() {
				var currentPath = $location.path();

				for (var i in $scope.menu) {
					var options = $scope.menu[i];

					for (var x in options) {
						var option = options[x];

						if (option.link === currentPath) {
							return option;
						}
					}
				}

				return null;
			}

			function generateTopMenu() {
				$scope.topMenu = $scope.menu[0].concat($scope.menu[1]);
			}

			function onLoad() {
				var option = determineActiveOption();
				setActiveOption(option);

				generateTopMenu();
			}

			onLoad();
		}
	]);
})(angular);
