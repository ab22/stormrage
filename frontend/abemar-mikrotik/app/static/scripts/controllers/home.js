; (function (angular) {
    'use strict';

    angular.module('app.controllers').controller('HomeCtrl', ['$scope', 'Mikrotik',
        function ($scope, Mikrotik) {
            $scope.clients = [];

            function requestClients() {
                Mikrotik.getClients().success(function (response) {
                    $scope.clients = response;
                });
            }

            $scope.reloadClients = function () {
                requestClients();
            };

            requestClients();
        }
    ]);
})(angular);
