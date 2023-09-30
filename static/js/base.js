var url = window.location.href;

//set active link class 
var els = document.querySelectorAll(".dropdown-menu a");
for (var i = 0, l = els.length; i < l; i++) {
    let el = els[i];
    if (el.href === url) {
       el.classList.add("active");
       let parent = el.closest(".main-nav"); // add this class for the top level "li" to get easy the parent
       if (parent){
        parent.classList.add("active");
       }
    }
}





// bootstrap snackBar 
var toastElList = [].slice.call(document.querySelectorAll('.toast'))
var toastList = toastElList.map(function (toastEl) {
return new bootstrap.Toast(toastEl, {
    'animation':true,
    'autohide':true,
    // 'delay': 5000, // delay sets by data attr
    })
})

// show bootstrap toast error 
function toastError(message){
    let toastEl = document.getElementById('errToastEl')
    if(toastEl){
        toastEl.classList.remove('bg-success', 'bg-primary', 'bg-danger');
        toastEl.classList.add('bg-danger');  
        toastEl.querySelector('.toast-body').innerHTML = message; 
        let toast = bootstrap.Toast.getInstance(toastEl)
        toast.show()
        return
    }
    console.error('Cant show error, element '+ elementID+' not found')
}

// show bootstrap toast info 
function toastInfo(message){
    let toastEl = document.getElementById('errToastEl')
    if(toastEl){
        toastEl.classList.remove('bg-success', 'bg-primary', 'bg-danger');
        toastEl.classList.add('bg-primary');  
        toastEl.querySelector('.toast-body').innerHTML = message; 
        let toast = bootstrap.Toast.getInstance(toastEl)
        toast.show()
        return
    }
    console.error('Cant show error, element '+ elementID+' not found')
}

// show bootstrap toast success 
function toastSuccess(message){
    let toastEl = document.getElementById('errToastEl')
    if(toastEl){
        toastEl.classList.remove('bg-success', 'bg-primary', 'bg-danger');
        toastEl.classList.add('bg-success');  
        toastEl.querySelector('.toast-body').innerHTML = message; 
        let toast = bootstrap.Toast.getInstance(toastEl)
        toast.show()
        return
    }
    console.error('Cant show error, element '+ elementID+' not found')
}