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
        var filename = file.replace(/^.*[\\\/]/, '')
        lastItem = null
        lastItem = file;


        console.log(file)
        $scope.showCard = false;
        angular.element(document.querySelector("#filename")).html("File: " + filename)



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
    //var output = [],
    //  i;



    container.html("")
    socket.onopen = function() {
        var filename = file.replace(/^.*[\\\/]/, '')
        container.append("<p><b>Tailing file: " + filename + "</b></p>");
        strf = file
        if (strf.indexOf("undefined") != 0) {
            container.append("<table class=\"fixed\" cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%'> " +
                "<col width=\"150px\" />" +
                "<col width=\"350px\" />" +
                "<col width=\"110px\" />" +
                "<col width=\"130px\" />" +
                "<col width=\"110px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"500px\" />" +
                "<col width=\"200px\" />" +
                "<tr > <td>" +
                "APPNAME" + "</td> <td>" +
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
        str = e.data.trim();
        if (str.indexOf("INFO") == 0) {

            str = str.replace("INFO", "");
            // str = "<table  cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%' >" +f2f3f4
            str = "<table  class=\"fixed\" cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%' bgcolor=\"#b0ffb0\" >" +
                /* "<col width=\"70%\" />" +
                "<col width=\"70%\" />" +
                "<col width=\"35%\" />" +
                "<col width=\"35%\" />" +
                "<col width=\"70%\" />" +
                "<col width=\"100%\" />" +
                "<col width=\"100%\" />" +
                "<col width=\"0%\" />" +
                "<col width=\"200px\" />" + */
                "<col width=\"150px\" />" +
                "<col width=\"350px\" />" +
                "<col width=\"50px\" />" +
                "<col width=\"130px\" />" +
                "<col width=\"100px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"500px\" />" +
                "<col width=\"200px\" />" +
                "<tr >" +
                str.replace(/,\n/g, "<tr >")
                .replace(/,/g, "<td width=\"100\" height=\"100\">")
                .replace(/<tr>$/, "") +
                "</table>";
            container.append(str);

        } else if (str.indexOf("ERROR") == 0) {
            str = str.replace("ERROR", "");
            str = "<table class=\"fixed\" cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%' bgcolor=\"#ffb0b0\" >" +
                "<col width=\"150px\" />" +
                "<col width=\"350px\" />" +
                "<col width=\"50px\" />" +
                "<col width=\"130px\" />" +
                "<col width=\"100px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"500px\" />" +
                "<col width=\"200px\" />" +
                "<tr >" +
                str.replace(/,\n/g, "<tr >")
                .replace(/,/g, "<td width=\"100\" height=\"100\">")
                .replace(/<tr>$/, "") +
                "</table>";
            container.append(str);

        } else if (str.indexOf("WARNING") == 0) {
            str = str.replace("WARNING", "");
            // str = "<table  cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%' bgcolor=\"#ffcc00\" >" +
            str = "<table class=\"fixed\" cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%' bgcolor=\"#ffff90\" >" +
                "<col width=\"150px\" />" +
                "<col width=\"350px\" />" +
                "<col width=\"50px\" />" +
                "<col width=\"130px\" />" +
                "<col width=\"100px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"500px\" />" +
                "<col width=\"200px\" />" +
                "<tr >" +
                str.replace(/,\n/g, "<tr >")
                .replace(/,/g, "<td width=\"100\" height=\"100\">")
                .replace(/<tr>$/, "") +
                "</table>";
            container.append(str);
            // container.append("<p style='background-color: yellow; color:blue'>" + str + "</p>" + "<hr>");
        } else if (str.indexOf("DEBUG") == 0) {
            str = str.replace("DEBUG", "");
            // str = "<table cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%' bgcolor=\"#00ff00\" >" +
            str = "<table class=\"fixed\" cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%' bgcolor=\"#a0a0a0\" >" +
                "<col width=\"150px\" />" +
                "<col width=\"350px\" />" +
                "<col width=\"50px\" />" +
                "<col width=\"130px\" />" +
                "<col width=\"100px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"500px\" />" +
                "<col width=\"200px\" />" +
                "<tr >" +
                str.replace(/,\n/g, "<tr >")
                .replace(/,/g, "<td width=\"100\" height=\"100\">")
                .replace(/<tr>$/, "") +
                "</table>";
            container.append(str);
        } else {
            container.append("<p style='background-color: yellow; color:blue'>" + str + "</p>" + "<hr>");

        }

        //container.append(str + "<br>" + "<hr>");

    }
    socket.onclose = function() {
        container.append("<p style='background-color: maroon; color:orange'>Connection Closed to WebSocket, tail stopped</p>");
    }
    socket.onerror = function(e) {
        container.append("<b style='color:red'>Some error occurred " + e.data.trim() + "<b>");
    }

    return socket;
}