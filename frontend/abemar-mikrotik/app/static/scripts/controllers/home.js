; (function(angular) {
    'use strict';

    angular.module('app.controllers').controller('HomeCtrl', ['$scope', 'Mikrotik', 'StrUtils',
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
                    }

                    $scope.clients = response.map(function(client) {
                        client.target = client.target.split(',');
                        client.max_limit = formatBytes(client.max_limit);
                        client.burst_limit = formatBytes(client.burst_limit);
                        client.burst_threshold = formatBytes(client.burst_threshold);

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
