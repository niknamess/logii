angular.module("logi2").controller("mainController", mainController);

mainController.$inject = ["$rootScope", "$scope", "$mdSidenav", "$http"]
const buttonR = document.getElementById('res');
const buttonErr = document.getElementById('btnerr');
const buttonInf = document.getElementById('btninf');
const buttonDbgs = document.getElementById('btndbgs');
const buttonWar = document.getElementById('btnwar');
const buttonView = document.getElementById('view');

var lastItem;
//const input = document.querySelector('input');

buttonR.addEventListener('click', event => {
    setTimeout(
        () => {
            window.location.reload();
        },
        1 * 200
    );
});

buttonView.addEventListener('click', event => {
    setTimeout(
        () => {
            //window.location.reload();
            initWS(lastItem)
            setBackColor('view', "#ed6c27")
            quotation('view', "Buttom")
        },
        1 * 200
    );
});

buttonErr.addEventListener('click', event => {
    setTimeout(
        () => {
            initWSType(lastItem, "ERROR", "#ffb0b0")
            setBackColor('view', "#ffb0b0")
            quotation('view', "ERROR")
        },
        1 * 200
    );
});
buttonInf.addEventListener('click', event => {
    setTimeout(
        () => {
            initWSType(lastItem, "INFO", "#b0ffb0")
            setBackColor('view', "#b0ffb0")
            quotation('view', "INFO")
        },
        1 * 200
    );
});

buttonDbgs.addEventListener('click', event => {
    setTimeout(
        () => {
            initWSType(lastItem, "DEBUG", "#a0a0a0")
            setBackColor('view', "#a0a0a0")
            quotation('view', "DEBUG")
        },
        1 * 200
    );
});

buttonWar.addEventListener('click', event => {
    setTimeout(
        () => {
            initWSType(lastItem, "WARNING", "#ffff90")
            setBackColor('view', "#ffff90")
            quotation('view', "WARNING")
                //setFontColor('view ', "black")
        },
        1 * 200
    );
});


function setBackColor(btn, color) {
    var property = document.getElementById(btn);
    property.style.backgroundColor = color

}

function setFontColor(btn, color) {
    var property = document.getElementById(btn);
    property.style.color = color

}

function quotation(id, text) {
    var q = document.getElementById(id);
    if (q) q.innerHTML = text;
}


function change(identifier, color) {
    identifier.style.background = color;
}

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

    var ws_proto = "ws:"
    if (window.location.protocol === "https:") {
        ws_proto = "wss:"
    }

    var socket = new WebSocket(ws_proto + "//" + window.location.hostname + ":" + window.location.port + "/ws/" + btoa(file));
    var container = angular.element(document.querySelector("#container"));




    container.html("")
    socket.onopen = function() {
        var filename = file.replace(/^.*[\\\/]/, '')
        container.append("<p><b>Tailing file: " + filename + "</b></p>");
        strf = file
        if (strf.indexOf("undefined") != 0) {
            container.append("<table class=\"fixed\" cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%'> " +
                "<col width=\"150px\" />" +
                "<col width=\"150px\" />" +
                "<col width=\"350px\" />" +
                "<col width=\"50px\" />" +
                "<col width=\"130px\" />" +
                "<col width=\"100px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"400px\" />" +
                "<col width=\"500px\" />" +
                "<col width=\"200px\" />" +
                "<tr > <td>" +
                "TYPE" + "</td> <td>" +
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

            str = str.replace("INFO", ",INFO");
            // str = "<table  cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%' >" +f2f3f4
            str = "<table  class=\"fixed\" cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%' bgcolor=\"#b0ffb0\" >" +
                "<col width=\"150px\" />" +
                "<col width=\"150px\" />" +
                "<col width=\"350px\" />" +
                "<col width=\"50px\" />" +
                "<col width=\"130px\" />" +
                "<col width=\"100px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"400px\" />" +
                "<col width=\"500px\" />" +
                "<col width=\"200px\" />" +
                "<tr >" +
                str.replace(/,\n/g, "<tr >")
                .replace(/,/g, "<td width=\"100\" height=\"100\">")
                .replace(/<tr>$/, "") +
                "</table>";
            container.append(str);

        } else if (str.indexOf("ERROR") == 0) {
            str = str.replace("ERROR", ",ERROR");
            str = "<table class=\"fixed\" cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%' bgcolor=\"#ffb0b0\" >" +
                "<col width=\"150px\" />" +
                "<col width=\"150px\" />" +
                "<col width=\"350px\" />" +
                "<col width=\"50px\" />" +
                "<col width=\"130px\" />" +
                "<col width=\"100px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"400px\" />" +
                "<col width=\"500px\" />" +
                "<col width=\"200px\" />" +
                "<tr >" +
                str.replace(/,\n/g, "<tr >")
                .replace(/,/g, "<td width=\"100\" height=\"100\">")
                .replace(/<tr>$/, "") +
                "</table>";
            container.append(str);

        } else if (str.indexOf("WARNING") == 0) {
            str = str.replace("WARNING", ",WARNING");
            // str = "<table  cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%' bgcolor=\"#ffcc00\" >" +
            str = "<table class=\"fixed\" cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%' bgcolor=\"#ffff90\" >" +
                "<col width=\"150px\" />" +
                "<col width=\"150px\" />" +
                "<col width=\"350px\" />" +
                "<col width=\"50px\" />" +
                "<col width=\"130px\" />" +
                "<col width=\"100px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"400px\" />" +
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
            str = str.replace("DEBUG", ",DEBUG");
            // str = "<table cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%' bgcolor=\"#00ff00\" >" +
            str = "<table class=\"fixed\" cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%' bgcolor=\"#a0a0a0\" >" +
                "<col width=\"150px\" />" +
                "<col width=\"150px\" />" +
                "<col width=\"350px\" />" +
                "<col width=\"50px\" />" +
                "<col width=\"130px\" />" +
                "<col width=\"100px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"400px\" />" +
                "<col width=\"500px\" />" +
                "<col width=\"200px\" />" +
                "<tr >" +
                str.replace(/,\n/g, "<tr >")
                .replace(/,/g, "<td width=\"100\" height=\"100\">")
                .replace(/<tr>$/, "") +
                "</table>";
            container.append(str);
        } else {
            container.append("<p style='background-color: #ffff90; color:blue'>" + str + "</p>" + "<hr>");

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







function initWSType(file, type, color) {

    var ws_proto = "ws:"
    if (window.location.protocol === "https:") {
        ws_proto = "wss:"
    }

    var socket = new WebSocket(ws_proto + "//" + window.location.hostname + ":" + window.location.port + "/ws/" + btoa(file));
    var container = angular.element(document.querySelector("#container"));




    container.html("")
    socket.onopen = function() {
        var filename = file.replace(/^.*[\\\/]/, '')
        container.append("<p><b>Tailing file: " + filename + "</b></p>");
        strf = file
        if (strf.indexOf("undefined") != 0) {
            container.append("<table class=\"fixed\" cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%'> " +
                "<col width=\"150px\" />" +
                "<col width=\"150px\" />" +
                "<col width=\"350px\" />" +
                "<col width=\"110px\" />" +
                "<col width=\"130px\" />" +
                "<col width=\"110px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"400px\" />" +
                "<col width=\"500px\" />" +
                "<col width=\"200px\" />" +
                "<tr > <td>" +
                "TYPE" + "</td> <td>" +
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
        if (str.indexOf(type) == 0) {

            str = str.replace(type, "");
            // str = "<table  cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%' >" +f2f3f4
            str = "<table  class=\"fixed\" cellspacing=\"0\" cellpadding=\"4\" border=\"1\" style='font-family:\"Courier New\", Courier, monospace; font-size:100%'" + "bgcolor=" + color + " >" +
                "<col width=\"150px\" />" +
                "<col width=\"150px\" />" +
                "<col width=\"350px\" />" +
                "<col width=\"50px\" />" +
                "<col width=\"130px\" />" +
                "<col width=\"100px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"400px\" />" +
                "<col width=\"500px\" />" +
                "<col width=\"200px\" />" +
                "<tr >" +
                str.replace(/,\n/g, "<tr >")
                .replace(/,/g, "<td width=\"100\" height=\"100\">")
                .replace(/<tr>$/, "") +
                "</table>";
            container.append(str);

        } else {
            //container.append("<p>" + "</p>");

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