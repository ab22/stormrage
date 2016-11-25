; (function (angular) {
	'use strict';

	angular.module('app.controllers').controller('MainLayoutCtrl', ['$scope', '$location', 'Auth',
		function ($scope, $location, Auth) {
			$scope.window = {};
			$scope.isCollapsed = true;

			function setActiveOption(option) {
				$scope.activeOption = option;
			}

			function generateTopMenu() {
				$scope.topMenu = $scope.menu[0].concat($scope.menu[1]);
			}

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

			function onLoad() {
				var option = determineActiveOption();
				setActiveOption(option);

				generateTopMenu();
			}

			$scope.signOut = function () {
				Auth.logout().success(function () {
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
						responsiveOnly: true,
					},
					{
						label: 'IP Privadas',
						icon: 'fa-users',
						link: '/main/privadas',
						responsiveOnly: true,
					},
					{
						label: 'Ping',
						icon: 'fa-users',
						link: '/main/ping',
						responsiveOnly: true,
					},
					{
						label: 'Clientes',
						icon: 'fa-user',
						link: '/main/home',
						responsiveOnly: true,
					},
					{
						label: 'Bloqueados',
						icon: 'fa-user-times',
						link: '/main/home',
						responsiveOnly: true,
					},
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
					$scope.isCollapsed = true;
				}
			}

			$scope.optionOnClick = function (option) {
				hideResponsiveMenu();

				if (typeof option.onClick !== 'undefined') {
					option.onClick();
					return;
				}

				setActiveOption(option);
			};

			$scope.isResponsiveMode = function () {
				return $scope.window.width <= 767;
			};

			$scope.showResponsiveOption = function (option) {
				var showResponsive = !option.responsiveOnly || (option.responsiveOnly && $scope.isResponsiveMode());

				return showResponsive;
			};

			onLoad();
		}
	]);
})(angular);
