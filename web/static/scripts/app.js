var app = angular.module('dashboard', ['ngMaterial']);

app.controller("dashboard", function($scope, $compile) {
    $scope.lines = [];
    console.log($scope);

    $scope.add_message = function(data) {
        $scope.lines.push({user: data['user'], msg: data['text']});
    }

    $scope.send_quit = function() {
        console.log("QUIT");
        var payload = {
            'event': 'quit'
        }
        socket.send(JSON.stringify(payload));
    };

    $scope.send_timeout = function(user) {
        console.log("timeout " + user);
        var payload = {
            'event': 'timeout',
            'data': {
                'target_user': user,
                'timeout_duration': '3'
            }

        }
        socket.send(JSON.stringify(payload));
    }

    connect_to_ws($scope);
});
