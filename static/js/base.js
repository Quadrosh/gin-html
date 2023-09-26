var url = window.location.href;

var els = document.querySelectorAll(".dropdown-menu a");
for (var i = 0, l = els.length; i < l; i++) {
    var el = els[i];
    if (el.href === url) {
       el.classList.add("active");
       var parent = el.closest(".main-nav"); // add this class for the top level "li" to get easy the parent
       parent.classList.add("active");
    }
}


var snackBarInfo = function (text) {
    var message = SnackBar({
        message: text,
        timeout: 10000,
        fixed: true,
    })
}
var snackBarError = function (text) {
    var message = SnackBar({
        status: "error",
        message: text,
        timeout: 30000,
        fixed: true,
    })
}