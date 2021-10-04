angular.module("logi2").controller("mainController", mainController);

mainController.$inject = ["$rootScope", "$scope", "$mdSidenav", "$http"]
const button = document.querySelector('button');
const input = document.querySelector('input');



function mainController($rootScope, $scope, $mdSidenav, $http) {
    var vm = this;

    vm.toggleSideNav = function toggleSideNav() {
        $mdSidenav('left').toggle()
    }



    vm.init = function init() {
        console.log("In the main controller")
        $scope.showCard = true;
        $http.get('searchproject')
            .then(function(result) {
                $rootScope.search_string = result.data["search_string"]
                console.log("Search :", result.data)
            }, function(result) {
                console.log("Failed to get search")
            })
    }

    vm.fontSize = ["10px", "11px", "12px", "14px", "16px", "18px", "20px", "22px", "24px"]
    $scope.currSize = vm.fontSize[2];

    $scope.open_connection = function(file) {
        console.log(file)
        $scope.showCard = false;
        // $scope.$apply()

        angular.element(document.querySelector("#filename")).html("File: " + file)
        var container = angular.element(document.querySelector("#container"))
        var ws;
        if (window.WebSocket === undefined) {
            container.append("Your browser does not support WebSockets");
            return;
        }

        input.addEventListener('click', event => {
                ws = initWS(file);
            }),


            input.addEventListener('change', event => {
                ws = initWS(file);
            }),
            //else {

            //  ws = initWS(file);

            //}
            input.addEventListener('input', event => {
                ws = initWS(file);
            });
        vm.toggleSideNav()
    }


    function updateValue() {
        ws = initWS(file);
    }

    function initWS(file) {

        var ws_proto = "ws:"
        if (window.location.protocol === "https:") {
            ws_proto = "wss:"
        }
        var socket = new WebSocket(ws_proto + "//" + window.location.hostname + ":" + window.location.port + "/ws/" + btoa(file));
        var container = angular.element(document.querySelector("#container"));

        // clear the contents
        container.html("")
        socket.onopen = function() {
            container.append("<p><b>Tailing file: " + file + "</b></p>");
        };
        socket.onmessage = function(e) {
            container.append(e.data.trim() + "<br>");
        }
        socket.onclose = function() {
            container.append("<p>Connection Closed to WebSocket, tail stopped</p>");
        }
        socket.onerror = function(e) {
            container.append("<b style='color:red'>Some error occurred " + e.data.trim() + "<b>");
        }
        return socket;
    }

    $scope.logout = function() {
        for (i = 0; i < document.forms.length; i++) {
            if (document.forms[i].id == "logoutForm") {
                document.forms[i].submit()
                return;
            }
        }
    }

    vm.init();
}