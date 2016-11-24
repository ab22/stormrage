; (function (angular) {
	'use strict';

	angular.module('app').config(['$stateProvider', '$urlRouterProvider',
		function ($stateProvider, $urlRouterProvider) {
			var viewsPath = 'static/views/';
			$urlRouterProvider.otherwise('/main/home');

			$stateProvider.
				state('main', {
					url: '/main',
					templateUrl: viewsPath + 'layout.html',
					controller: 'MainLayoutCtrl',
					requiresAuthentication: true
				}).
				state('main.home', {
					url: '/home',
					templateUrl: viewsPath + 'privates.html',
					controller: 'PrivatesCtrl',
					requiresAuthentication: true
				}).
				state('main.privates', {
					url: '/privates',
					templateUrl: viewsPath + 'privates.html',
					controller: 'PrivatesCtrl',
					requiresAuthentication: true
				}).state('login', {
					url: '/login',
					templateUrl: viewsPath + 'login.html',
					controller: 'LoginCtrl',
					requiresAuthentication: false
				});
		}
	]);
})(angular);
