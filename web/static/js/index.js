angular.module("logi2").controller("mainController", mainController);

mainController.$inject = ["$rootScope", "$scope", "$mdSidenav", "$http"]
const button = document.getElementById('btn');
const buttonR = document.getElementById('res');

var lastItem;
//const input = document.querySelector('input');


button.addEventListener('click', event => {
    setTimeout(
        () => {
            initWS(lastItem)
        },
        1 * 200
    );
});

buttonR.addEventListener('click', event => {
    setTimeout(
        () => {
            window.location.reload();
        },
        1 * 200
    );
});

function mainController($rootScope, $scope, $mdSidenav, $http) {

    var vm = this;
    //var lastItem;

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
        lastItem = null
        lastItem = file;


        console.log(file)
        $scope.showCard = false;
        angular.element(document.querySelector("#filename")).html("File: " + file)



        var container = angular.element(document.querySelector("#container"))

        var ws;
        if (window.WebSocket === undefined) {
            container.append("Your browser does not support WebSockets");
            return;
        } else {


            ws = initWS(file);

        }



        vm.toggleSideNav()
    }

    vm.init();
}

function initWS(file) {
    // document.querySelector('button').removeEventListener(initWS(file));

    //window.alert("InitWs Files" + file);
    var ws_proto = "ws:"
    if (window.location.protocol === "https:") {
        ws_proto = "wss:"
    }

    var socket = new WebSocket(ws_proto + "//" + window.location.hostname + ":" + window.location.port + "/ws/" + btoa(file));
    var container = angular.element(document.querySelector("#container"));

    container.html("")
    socket.onopen = function() {
        container.append("<p><b>Tailing file: " + file + "</b></p>");
        strf = file
        if (strf.indexOf("undefined") != 0) {
            container.append("<table border=\"1\"> <tr > <td>" +
                "TYPE MESSAGE" + "</td> <td >" +
                "APPNAME" + "</td> <td >" +
                "APPPATH" + "</td> <td>" +
                "APPPID" + "</td><td>" +
                "THREAD" + "</td><td>" +
                "TIME" + "</td><td>" +
                "ULID" + "</td><td>" +
                "MESSAGE" + "</td><td>" +
                "DETAILS" + "</td></tr > </table >");
        }
    }

    socket.onmessage = function(e) {
        //  let msg = e.data.trim();
        str = e.data.trim();
        if (str.indexOf("INFO") == 0) {


            //str.css("background-color", 'red');
            container.append("<p style='background-color: white; color:black'>" + str + "</p>" + "<hr>");

        } else if (str.indexOf("ERROR") == 0) {
            //str.css("background-color", 'orange');
            container.append("<p style='background-color: maroon; color:orange'>" + str + "</p>" + "<hr>");
            //}


        } else if (str.indexOf("WARNING") == 0) {
            //str.css("background-color", 'yellow');
            container.append("<p style='background-color: yellow; color:blue'>" + str + "</p>" + "<hr>");
        } else {
            container.append(str + "<hr>");

        }

        //container.append(str + "<br>" + "<hr>");

    }
    socket.onclose = function() {
        container.append("<p style='background-color: maroon; color:orange'>Connection Closed to WebSocket, tail stopped</p>");
    }
    socket.onerror = function(e) {
        container.append("<b style='color:red'>Some error occurred " + e.data.trim() + "<b>");
    }


    //  window.alert("Socket " + socket);
    return socket;


}

d3.text(str, function(data) {
    var parsedCSV = d3.csv.parseRows(data);

    var container = d3.select("container")
        .append("table")

    .selectAll("tr")
        .data(parsedCSV).enter()
        .append("tr")

    .selectAll("td")
        .data(function(d) { return d; }).enter()
        .append("td")
        .text(function(d) { return d; });
});