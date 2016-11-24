; (function(angular) {
    'use strict';

    angular.module('app.controllers').controller('HomeCtrl', ['$scope', 'Mikrotik',
        function($scope, Mikrotik) {
            $scope.clients = [];
            $scope.rowCollection = [];

            function requestClients() {
                Mikrotik.getClients().success(function(response) {
                    response.forEach(function(element) {

                    }, this);

                    $scope.clients = response.map(function(client) {
                        client.target = client.target.split(',');
                        $scope.rowCollection.push(client);

                        return client;
                    });

                });
            }

            $scope.reloadClients = function() {
                requestClients();
            };

            requestClients();
        }
    ]);
})(angular);
