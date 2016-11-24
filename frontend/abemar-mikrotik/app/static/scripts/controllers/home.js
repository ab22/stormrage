; (function (angular) {
    'use strict';

    angular.module('app.controllers').controller('HomeCtrl', ['$scope', 'Mikrotik',
        function ($scope, Mikrotik) {
            $scope.clients = [];

            function requestClients() {
                Mikrotik.getClients().success(function (response) {
                    response.forEach(function (element) {

                    }, this);

                    $scope.clients = response.map(function (client) {
                        client.target = client.target.split(',')
                        return client;
                    });
                });
            }

            $scope.reloadClients = function () {
                requestClients();
            };

            requestClients();
        }
    ]);
})(angular);
