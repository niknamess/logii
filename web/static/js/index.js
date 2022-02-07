angular.module("logi2").controller("mainController", mainController);

mainController.$inject = ["$rootScope", "$scope", "$mdSidenav", "$http"]
const buttonR = document.getElementById('res');
const buttonErr = document.getElementById('btnerr');
const buttonInf = document.getElementById('btninf');
const buttonDbgs = document.getElementById('btndbgs');
const buttonWar = document.getElementById('btnwar');
const buttonAll = document.getElementById('btnall');
const buttonClr = document.getElementById('changeclr');
const inputform = document.getElementById('search_string');
const textList = document.getElementById('listfile');
//const trInf = document.getElementById('trInf');
var countWar = 0
var countErr = 0
var countInf = 0
var countDbg = 0
var countFtl = 0
var countAll = 0
var start
var standartform = ""
var lastItem;
var statusS = "empty"
    //const input = document.querySelector('input');



function isEmpty(str) {
    return (!str || 0 === str.length);
}

buttonR.addEventListener('click', event => {
    setTimeout(
        () => {
            window.location.reload();

        },
        1 * 200
    );
});


inputform.addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        // code for enter
        setTimeout(
            () => {
                Null()
                initWS(lastItem, statusS)
                Null()
            },
            1 * 200
        );
    }
});
/* 
buttonErr.addEventListener('click', event => {

    setTimeout(
        () => {
            Null()
            initWS(lastItem, "ERROR")
            statusS = "ERROR"
            setBackColor('changeclr', "#ffb0b0")

        },
        1 * 200
    );
});
buttonInf.addEventListener('click', event => {
    setTimeout(
        () => {
            Null()
            initWS(lastItem, "INFO")
            statusS = "INFO"
            setBackColor('changeclr', "#b0ffb0")

        },
        1 * 200
    );
});

buttonDbgs.addEventListener('click', event => {
    setTimeout(
        () => {
            Null()
            statusS = "DEBUG"
            initWS(lastItem, "DEBUG")
            setBackColor('changeclr', "#a0a0a0")


        },
        1 * 200
    );
});

buttonWar.addEventListener('click', event => {
    setTimeout(
        () => {
            Null()
            statusS = "WARNING"
            initWS(lastItem, "WARNING")
            setBackColor('changeclr', "#ffff90")


        },
        1 * 200
    );
});
buttonAll.addEventListener('click', event => {
    setTimeout(
        () => {
            Null()
            statusS = "empty"
            initWS(lastItem, "empty")
            setBackColor('changeclr', "#ed6c27")


        },
        1 * 200
    );
});
 */
/* trInf.addEventListener('click', event => {
    setTimeout(
        () => {
            Null()
            initWS(lastItem, "INFO")
            statusS = "INFO"

        },
        1 * 200
    );
}); */
function editInf() {
    Null()
    initWS(lastItem, "INFO")
    statusS = "INFO"
}

function editErr() {
    Null()
    initWS(lastItem, "ERROR")
    statusS = "ERROR"
}

function editDbgs() {
    Null()
    initWS(lastItem, "DEBUG")
    statusS = "DEBUG"
}

function editWarn() {
    Null()
    initWS(lastItem, "WARNING")
    statusS = "WARNING"
}

function editAll() {
    Null()
    initWS(lastItem, "empty")
    statusS = "empty"
}

function myFunction() {
    document.getElementById("trInf").innerHTML = "YOU CLICKED ME!";
}

function Null() {
    countWar = 0
    countErr = 0
    countInf = 0
    countDbg = 0
    countFtl = 0
    countAll = 0
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

    // vm.fontSize = ["10px", "11px", "12px", "14px", "16px", "18px", "20px", "22px", "24px"]
    // $scope.currSize = vm.fontSize[2];


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
            ws = initWS(file, "empty");

        }

        vm.toggleSideNav()
    }

    vm.init();
}

function initWS(file, type) {

    var observer = new MutationObserver(function(_mutations, me) {
        // `mutations` is an array of mutations that occurred
        // `me` is the MutationObserver instance
        start = document.getElementById('Foxtrot');
        if (start) {
            handleCanvas(start);
            me.disconnect(); // stop observing
            return;
        }
    });

    // start observing
    observer.observe(document, {
        childList: true,
        subtree: true
    });


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

            container.append("<table");

        }
    }

    socket.onmessage = function(e) {
        var loglist
        str = e.data.trim();


        parser = new DOMParser();
        xmlDoc = parser.parseFromString(str, "text/xml");
        loglist = xmlDoc.getElementsByTagName("loglist")

        // document.getElementById('follow').scrollIntoView();
        k2 = isEmpty(loglist)
        if (k2 == false) {
            Null()
            str = ParseXml(str, type)
                /*  if (type == "INFO" || type == "empty") {
                     countInf = countInf / 2
                 } */
                //document.getElementById("clear1").innerHTML = "";
            countWar = countWar / 2
            countErr = countErr / 2
            countInf = countInf / 2
            countDbg = countDbg / 2
            countFtl = countFtl / 2
            countAll = countAll / 2
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
                "<tr > <td class=\"info\" onclick=\"editInf()\">" + "INFO:" +
                countInf + "</td> <td class=\"error\" onclick=\"editErr()\">" + "Error:" +
                countErr + "</td> <td class=\"warning\" onclick=\"editWarn()\">" + "Warning:" +
                countWar + "</td> <td class=\"debug\" onclick=\"editDbgs()\">" + "Debug:" +
                countDbg + "</td> <td class=\"all\" onclick=\"editAll()\">" + "All:" +
                countAll +
                "</td></tr > </table >");

            Null()
            container.append("<div style=\"\" class=\"TableContainer\">" +
                "<table id=\"tbl92\" border=\"0\" class=\"tableScroll\"  data-scroll-speed=2 align=\"center\" >" +
                "<thead>" +
                "<tr>" +
                "<th onclick=\"Vi.Table.sort.string(this)\" title=\"Strings will be ordered lessically.\" > TYPE </th>" +
                "<th onclick=\"Vi.Table.sort.string(this)\" title=\"Strings will be ordered lessically.\" > APPNAME </th>" +
                "<th onclick=\"Vi.Table.sort.string(this)\" title=\"Strings will be ordered lessically.\" > APPPATH </th>" +
                "<th onclick=\"Vi.Table.sort.number(this)\" title=\"Number will be sortes as number.\" > APPPID </th>" +
                "<th class = \"th-sm\" > THREAD </th>" +
                "<th class = \"th-sm\" > TIME </th>" +
                "<th class = \"th-sm\" > ULID </th>" +
                "<th class = \"th-sm\" > MESSAGE </th>" +
                "<th class = \"th-sm\" > DETAILS </th> </tr>" +
                "</thead>" +
                "<tbody>" + str + "</tbody></table></div>");
            container.append(
                "<div class =\"sysinfo d-block p-2 .bg-light.bg-gradient text-dark\">" + "Message: " +
                "</div>" +
                "<div class =\"sysinfo d-block p-2 .bg-light.bg-gradient text-dark\" id = \"message\">" +
                "</div>" +
                "<div class =\"sysinfo d-block p-2 .bg-secondary.bg-gradient text-dark\" id = \"details\">" +
                "</div>");


        } else {
            if (str == "Indexing file, please wait") {
                container.append(" <div id =\"load\" class=\"center\">" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"textL\">Loading...</div>" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"wave\"></div>" +
                    "</div>");

            } else if (str == "Indexing complated") {

                document.getElementById("load").remove();
                container.append("<div class=\"textL\">Indexing complated!</div>");
            } else
                container.append("<hr>" + str + "</hr>");
        }

        container.append(standartform);

    }
    socket.onclose = function() {
        container.append("<p style='background-color: maroon; color:orange'>Connection Closed to WebSocket, tail stopped</p>");
    }
    socket.onerror = function(e) {
        container.append("<b style='color:red'>Some error occurred " + e.data.trim() + "<b>");
    }

    return socket;
}



function ParseXml(str, type) {
    var parser, xmlDoc, table, heyho;
    parser = new DOMParser();
    xmlDoc = parser.parseFromString(str, "application/xml");
    log = xmlDoc.getElementsByTagName("log");
    for (i = 0; i < log.length; i++) {
        if (i == 0) {
            heyho = "id='Foxtrot'"
        } else {
            heyho = ""
        }
        if (type == typeMsg(log[i].getAttribute('type')) || type == "empty") {
            table +=
                "<tr " + heyho + "  bgcolor =" + Color(log[i].getAttribute('type')) + ">" + "<td class=\"\"><span>" +
                typeMsg(log[i].getAttribute('type')) +
                "</span></td><td class=\"\"><span>" +
                log[i].getAttribute('module_name') +
                "</span></td><td class=\"\"><span>" +
                log[i].getAttribute('app_path') +
                "</span></td><td class=\"ellipsis\"><span>" +
                log[i].getAttribute('app_pid') +
                "</span></td><td class=\"\"><span>" +
                log[i].getAttribute('thread_id') +
                "</span></td><td class=\"\"><span>" +
                split_at_index(log[i].getAttribute('time')) +
                "</span></td><td class=\"\"><span>" +
                log[i].getAttribute('ulid') +
                "</span></td><td class=\"ellipsis\"><span>" +
                log[i].getAttribute('message') +
                "</span></td><td class=\"ellipsis\"><span>" +
                log[i].getAttribute('ext_message') +
                "</span></td></tr>";
        }
    }

    return table
}



//23072021005653.991
//(year, monthIndex, day, hours, minutes, seconds, milliseconds)
//23 07 2021 00 25 53.492
function split_at_index(value) {
    norm = value.substring(0, 2) + "." + value.substring(2);
    norm = norm.substring(0, 5) + "." + norm.substring(5);
    norm = norm.substring(0, 10) + " " + norm.substring(10);
    norm = norm.substring(0, 13) + ":" + norm.substring(13);
    norm = norm.substring(0, 16) + ":" + norm.substring(16);
    norm = norm.substring(0, 19) + ":" + norm.substring(19);
    norm = norm.substring(0, 23) + " UTC" + norm.substring(23);
    return norm
}

function typeMsg(type) {
    if (type == "1") {
        msg = "DEBUG";
        countDbg++
    } else if (type == "0") {
        msg = "INFO";
        countInf++
    } else if (type == "2") {
        msg = "WARNING";
        countWar++
    } else if (type == "3") {
        msg = "ERROR";
        countErr++
    } else if (type == "4") {
        msg = "FATAL";
        countFtl++
    }
    countAll++
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

        color = "#b2ffb2";
    }
    return color
}