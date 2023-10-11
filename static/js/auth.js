



function setCookie(name,value,days) {
    var expires = "";
    if (days) {
        var date = new Date();
        date.setTime(date.getTime() + (days*24*60*60*1000));
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = name + "=" + (value || "")  + expires + "; path=/";
}

function getCookie(name) {
    var nameEQ = name + "=";
    var ca = document.cookie.split(';');
    for(var i=0;i < ca.length;i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1,c.length);
        if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
    }
    return null;
}

// Save auth cookie, calls during handle signin response
function setAccessToken(token){
    setCookie('auth', token, 3)
}

function logout(){
    setCookie('auth', '', 1)
    window.location = window.location.href
}

// bootstrap form validation
(function () {
    'use strict'
    const forms = document.querySelectorAll('.requires-validation')
    Array.from(forms)?.forEach( form => {
        form.addEventListener('submit', function (event) {
            if (!form.checkValidity()) {
                event.preventDefault()
                event.stopPropagation()
            }
            form.classList.add('was-validated')
        }, false)
      })


})();

(function () {
    'use strict'
    const signOutButton = document.getElementById('signOutButton')
    const signInButton = document.getElementById('signInButton')
    const adminButton = document.getElementById('adminButton')
    let authCookie = getCookie('auth')
    if (authCookie.length){
        signOutButton?.classList.remove('d-none')
        signInButton?.classList.add('d-none')
        adminButton?.classList.remove('d-none')
    } else {
        signOutButton?.classList.add('d-none')
        signInButton?.classList.remove('d-none')
        adminButton?.classList.add('d-none')

    }
   
})();