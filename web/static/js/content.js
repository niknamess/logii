var start = document.getElementById('Foxtrot')
    // callback executed when canvas was found
function handleCanvas(start) {
    // start = document.getElementById('Foxtrot');
    start.focus();
    start.style.backgroundColor = 'green';
    document.onkeydown = checkKey;
}

// set up the mutation observer
var observer = new MutationObserver(function(mutations, me) {
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



function dotheneedful(sibling) {
    if (sibling != null) {
        start.focus();
        start.style.backgroundColor = '';
        start.style.color = '';
        sibling.focus();
        sibling.style.backgroundColor = 'green !important';
        start = sibling;
    }
}



function checkKey(e) {
    e = e || window.event;
    if (e.keyCode == '37') {
        // left arrow
        var sibling = start.previousElementSibling;
        dotheneedful(sibling);
    } else if (e.keyCode == '39') {
        // right arrow
        var sibling = start.nextElementSibling;
        dotheneedful(sibling);
    }
}