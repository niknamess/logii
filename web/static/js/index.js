angular.module("logi2").controller("mainController", mainController);

mainController.$inject = ["$rootScope", "$scope", "$mdSidenav", "$http"]
const button = document.querySelector('button');
//const input = document.querySelector('input');



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

    // button.addEventListener('click', event => {
    //    console.log(file)


    // });

    $scope.click_fucking_button = function(file) {

            // container.empty();
            //  socket.html("")

            window.alert("Button for click  " + file);
            ws = initWSE(file);
            window.alert("Button" + file);


        }
        //<button (click)=”handleClick($event)” type=”button” >Обновить</button> 
        //<button (click)=”handleClick($event)” ng-click=”button('<<.>>')” >Обновить</button> 
    $scope.open_connection = function(file) {


        console.log(file)
        $scope.showCard = false;
        // $scope.$apply()

        window.alert("Start");
        // initWSE(file) 

        //button.removeEventListener();

        button.addEventListener('click', event => {
            window.alert("Start event" + event);
            window.alert("File button" + file);
            ws = initWSE();
        });

        //angular.element(document.querySelector("#filename")).html("File: " + file)
        initWSE()

        angular.element(document.querySelector("#filename")).html("File: " + file)


        button.removeEventListener();

        function initWSE() {
            var container = angular.element(document.querySelector("#container"))
                // var socket = new WebSocket(ws_proto + "//" + window.location.hostname + ":" + window.location.port + "/ws/" + btoa(file));

            var ws;
            if (window.WebSocket === undefined) {
                container.append("Your browser does not support WebSockets");
                return;
            } else {


                ws = initWS(file);
                //button.addEventListener('click', event => {
                //   ws = initWS(file);
                //});
                //button.removeEventListener('click', event);
                //socket.empty();
                //container.empty();
            }



            //button.addEventListener('click', event => {
            // container.empty();
            //  socket.html("")

            //  window.alert("Button for click  " + file);
            //ws = initWS(file);
            //window.alert("Button" + file);
            //event.currentTarget.removeEventListener(event.type);

            //});

            //button.removeEventListener('click', event);
            // socket.html("")





            //  vm.toggleSideNav()
            // socket.empty();
            //container.empty();
            // container.empty();
        }
        vm.toggleSideNav()
    }



    function initWS(file) {
        window.alert("InitWs Files" + file);
        var ws_proto = "ws:"
        if (window.location.protocol === "https:") {
            ws_proto = "wss:"
        }
        //remove(file)
        var socket = new WebSocket(ws_proto + "//" + window.location.hostname + ":" + window.location.port + "/ws/" + btoa(file));
        var container = angular.element(document.querySelector("#container"));
        //socket.empty()
        //container.empty()
        //clear the contents
        container.html("")
        socket.onopen = function() {
            container.append("<p><b>Tailing file: " + file + "</b></p>");

        }
        socket.onmessage = function(e) {
            container.append(e.data.trim() + "<br>");
        }
        socket.onclose = function() {
            container.append("<p>Connection Closed to WebSocket, tail stopped</p>");
        }
        socket.onerror = function(e) {
                container.append("<b style='color:red'>Some error occurred " + e.data.trim() + "<b>");
            }
            //container.empty()
            //remove(file)

        window.alert("Socket " + socket);
        return socket;


    }

    // $scope.logout = function() {
    //   for (i = 0; i < document.forms.length; i++) {
    //     if (document.forms[i].id == "logoutForm") {
    //       document.forms[i].submit()
    //     return;
    //}
    //}
    //}

    vm.init();
}