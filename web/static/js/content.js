function handleCanvas(start) {
    start.focus();
    start.style.backgroundColor = 'green';
    document.getElementById("message").innerHTML = start.cells.item(7).textContent
    document.getElementById("details").innerHTML = start.cells.item(8).textContent
    document.onkeydown = checkKey;
}
var count = 0

function dotheneedful(sibling) {
    if (sibling != null) {
        start.focus();
        start.style.backgroundColor = '';
        start.style.color = '';
        start.id = "";
        sibling.focus();
        sibling.style.backgroundColor = 'green';
        sibling.style.color = 'white';
        sibling.id = "follow"
        sibling.style.backgroundColor = 'green';
        sibling.style.color = 'white';
        start = sibling;
        standartform = start

        document.getElementById("message").innerHTML = start.cells.item(7).textContent
        document.getElementById("details").innerHTML = start.cells.item(8).textContent
        var elmnt = document.getElementById("follow");
        elmnt.scrollIntoView(false);
    }

}

document.onkeydown = checkKey;
//37 39
function checkKey(e) {
    e = e || window.event;
    if (e.keyCode == '38') {
        // up arrow
        if (count != 0) {
            count = count - 1
        }
        var sibling = start.previousElementSibling;
        dotheneedful(sibling);
    } else if (e.keyCode == '40') {
        // down arrow
        count = count + 1
        var sibling = start.nextElementSibling;
        console.log(sibling)
        if (sibling != "null") {
            count = count + 1
        }
        dotheneedful(sibling);
        //elem.addEventListener('click', function hideContent(e)
    }
    /* else if (e.addEventListener('click')) {

           var sibling = start.ElementSibling;
           dotheneedful(sibling);
       } */
    console.log(count)
}


//add
// This is an example of a custom sort. Despite the values in this column are string,  
// the column will be ordered as they were numbers (instead of lessically)
function sortCustom1(th) {
    try {

        // the column will be ordered following the same order the items are in the array. 
        var numbers = ['zero', 'one', 'two', 'three', 'four', 'five', 'six', 'seven', 'eight', 'nine'];

        function getValue(tr, cellIndex) {
            var value = tr.children[cellIndex].innerText.toLowerCase().trim();
            return numbers.indexOf(value);
        }

        Vi.Table.sort(th, getValue);

    } catch (jse) {
        console.error(jse);
    }
}

/**
 * Here the column is sorted based on an atribute and not the value shown.
 * That should highlight the fact the developer has an hight degree of 
 * freedom on how implement the table and the data 
 * At the end, the only constrain is that the function 'getValue' must
 * return a sortable value.
 */
function sortCustom2(th) {
    try {

        function getValue(tr, cellIndex) {
            var child = tr.children[cellIndex];
            var ticks = child.getAttribute("data-ticks");
            var value = parseInt(ticks);
            return value;
        }

        Vi.Table.sort(th, getValue);

    } catch (jse) {
        console.error(jse);
    }
}