angular.module("logi2").controller("mainController", mainController);

mainController.$inject = ["$rootScope", "$scope", "$mdSidenav", "$http"]
const buttonR = document.getElementById('res');
const buttonErr = document.getElementById('btnerr');
const buttonInf = document.getElementById('btninf');
const buttonDbgs = document.getElementById('btndbgs');
const buttonWar = document.getElementById('btnwar');
const buttonView = document.getElementById('view');
var countWar = 0
var countErr = 0
var countInf = 0
var countDbg = 0

var lastItem;
//const input = document.querySelector('input');


function isEmpty(str) {
    return (!str || 0 === str.length);
}

buttonR.addEventListener('click', event => {
    setTimeout(
        () => {
            window.location.reload();
            Null()
        },
        1 * 200
    );
});

buttonView.addEventListener('click', event => {

    setTimeout(
        () => {

            Null()
            countWS(lastItem)
            initWS(lastItem)
            setBackColor('view', "#ed6c27")
            quotation('view', "Find")
        },
        1 * 200
    );




});

buttonErr.addEventListener('click', event => {
    Null()
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

function Null() {
    countWar = 0
    countErr = 0
    countInf = 0
    countDbg = 0
}

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
    Null()
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
            // Null()
            countWS(file)
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
            container.append("<table  > " +
                "<col width=\"150px\" />" +
                "<col width=\"150px\" />" +
                "<col width=\"350px\" />" +
                "<col width=\"100px\" />" +
                "<col width=\"130px\" />" +
                "<col width=\"100px\" />" +
                "<col width=\"300px\" />" +
                "<col width=\"400px\" />" +
                "<col width=\"500px\" />" +
                "<col width=\"200px\" />" +
                "<tr>" +
                "<th class = \"th-sm\" > TYPE </th>" +
                "<th class = \"th-sm\" > APPNAME </th>" +
                "<th class = \"th-sm\" > APPPATH </th>" +
                "<th class = \"th-sm\" > APPPID </th>" +
                "<th class = \"th-sm\" > THREAD </th>" +
                "<th class = \"th-sm\" > TIME </th>" +
                "<th class = \"th-sm\" > ULID </th>" +
                "<th class = \"th-sm\" > MESSAGE </th>" +
                "<th class = \"th-sm\" > DETAILS </th> </tr>" +
                "</table>");

        }
    }

    socket.onmessage = function(e) {
        var loglist
        str = e.data.trim();


        parser = new DOMParser();
        xmlDoc = parser.parseFromString(str, "text/xml");
        loglist = xmlDoc.getElementsByTagName("loglist")


        //k2 = isEmpty(loglist)
        //  if (k2 == false) {
        //str = ParseXml(str)
        //container.append(str);
        //} else {
        container.append("<br>" + str + "</br>" + "<hr>");
        //  }


    }
    socket.onclose = function() {
        container.append("<p style='background-color: maroon; color:orange'>Connection Closed to WebSocket, tail stopped</p>");
        Null()
    }
    socket.onerror = function(e) {
        container.append("<b style='color:red'>Some error occurred " + e.data.trim() + "<b>");
    }

    return socket;
}


function countWS(file) {

    var ws_proto = "ws:"
    if (window.location.protocol === "https:") {
        ws_proto = "wss:"
    }

    var socket = new WebSocket(ws_proto + "//" + window.location.hostname + ":" + window.location.port + "/ws/" + btoa(file));
    var container = angular.element(document.querySelector("#container"));




    container.html("")
    socket.onopen = function() {
        var filename = file.replace(/^.*[\\\/]/, '')
        strf = file
        if (strf.indexOf("undefined") != 0) {
            container.append("<table > " +
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
                "<tr > <td class=\"info\">" + "INFO:" +
                countInf + "</td> <td class=\"error\">" + "Error:" +
                countErr + "</td> <td class=\"warning\">" + "Warning:" +
                countWar + "</td> <td class=\"debug\">" + "Debug:" +
                countDbg +
                "</td></tr > </table >");
        }
    }


    socket.onclose = function() {
        container.append("<p style='background-color: maroon; color:orange'>" + "Connection Closed to WebSocket, tail stopped" + "</p>");
        Null()
    }
    socket.onerror = function(e) {
        container.append("<b>Some error occurred " + e.data.trim() + "<b>");
    }

    return socket;
}


/* function ParseXml(xml) {
    var parser, xmlDoc, table, msg;
    parser = new DOMParser();
    xmlDoc = parser.parseFromString(xml, "text/xml");
    loglist = xmlDoc.getElementsByTagName("loglist");
    for (i = 0; i < loglist.length; i++) {
        log = loglist[i].getElementsByTagName("log");
        for (i = 0; i < log.length; i++) {
            table +=
                "<tbody bgcolor =" + typeMsg(log[i].getAttribute('type')) + ">" +
                "<tr><td>" +
                typeMsg(log[i].getAttribute('type')) + //type
                "</td><td>" +
                log[i].getAttribute('module_name') + //module_name
                "</td><td>" +
                log[i].getAttribute('app_path') + //app_path
                "</td><td>" +
                log[i].getAttribute('app_pid') + //app_pid
                "</td><td>" +
                log[i].getAttribute('thread_id') + //thread_id
                "</td><td>" +
                log[i].getAttribute('time') +
                "</td><td>" +
                log[i].getAttribute('ulid') +
                "</td><td>" +
                log[i].getAttribute('message') + //message
                "</td><td>" +
                log[i].getAttribute('ext_message') + //ext_message
                "</td></tr>" +
                "</tbody>";

        }
    }

    return table

} */

function ParseXml(str) {

    //:TODO create sortable table with xml structure
    var parser, xmlDoc, table;
    parser = new DOMParser();
    //serializer = new XMLSerializer();
    xmlDoc = parser.parseFromString(str, "application/xml");
    //xmlDoc1 = serializer.serializeToString(xmlDoc);
    //loglist = xmlDoc.getElementsByTagName("loglist");
    log = xmlDoc.getElementsByTagName("log");
    for (i = 0; i < log.length; i++) {
        table += //< table > < /table>
            //"<tbody>" +
            "<table bgcolor =" + typeMsg(log[i].getAttribute('type')) + ">" +
            "<tr><td>" +
            typeMsg(log[i].getAttribute('type')) + //type
            "</td><td>" +
            log[i].getAttribute('module_name') + //module_name
            "</td><td>" +
            log[i].getAttribute('app_path') + //app_path
            "</td><td>" +
            log[i].getAttribute('app_pid') + //app_pid
            "</td><td>" +
            log[i].getAttribute('thread_id') + //thread_id
            "</td><td>" +
            log[i].getAttribute('time') +
            "</td><td>" +
            log[i].getAttribute('ulid') +
            "</td><td>" +
            log[i].getAttribute('message') + //message
            "</td><td>" +
            log[i].getAttribute('ext_message') +
            "</td><td>" +
            log[i].getAttribute('ddMMyyyyhhmmsszzz') + //message

            "</td></tr>" +
            "</table>";

    }


    return table

}


function typeMsg(type) {
    if (type == "0") {
        msg = "INFO";
        countInf++
    } else if (type == "1") {
        msg = "DEBUG";
        countDbg++
    } else if (type == "2") {
        msg = "WARNING";
        countWar++

    } else if (type == "3") {
        msg = "ERROR";
        countErr++
    } else if (type == "4") {
        msg = "FATAL";

    }
    return msg
}

function Color(type) {
    if (type == "0") {

        color = "#b4fcb5";
    } else if (type == "1") {

        color = "#a0a0a0";
    } else if (type == "2") {

        color = "#fffc9b";
    } else if (type == "3") {

        color = "#fdb1b1";
    } else if (type == "4") {

        color = "#b1ffb1";
    }
    return color
}