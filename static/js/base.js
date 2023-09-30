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







/* Example toast */
// (new ToastError('text')).show()
// (new ToastInfo('text')).show()



/* Toast constuctor */
function ToastInfo(text) {
    var toastElement = _buildToast(text, 'bg-primary', 10000); // 'bg-success', 'bg-primary', 'bg-danger'
    var toastWrapper = _getOrCreateToastWrapper();
    toastWrapper.append(toastElement);
    this.bootstrapToast = bootstrap.Toast.getOrCreateInstance(toastElement);
    this.show = function() {
        this.bootstrapToast.show();
    }
    this.hide = function() {
        this.bootstrapToast.hide();
    }
    this.dispose = function() {
        this.bootstrapToast.dispose();
    }
}

function ToastError(text) {
    var toastElement = _buildToast(text, 'bg-danger', 15000); // 'bg-success', 'bg-primary', 'bg-danger'
    var toastWrapper = _getOrCreateToastWrapper();
    toastWrapper.append(toastElement);
    this.bootstrapToast = bootstrap.Toast.getOrCreateInstance(toastElement);
    this.show = function() {
        this.bootstrapToast.show();
    }
    this.hide = function() {
        this.bootstrapToast.hide();
    }
    this.dispose = function() {
        this.bootstrapToast.dispose();
    }
}


/* Toast Utility methods */
function _getOrCreateToastWrapper() {
    var toastWrapper = document.querySelector('body > [data-toast-wrapper]');
    if (!toastWrapper) {
        toastWrapper = document.createElement('div');
        toastWrapper.style.zIndex = 11;
        toastWrapper.style.position = 'fixed';
        toastWrapper.style.bottom = 0;
        toastWrapper.style.right = 0;
        toastWrapper.style.padding = '1rem';
        toastWrapper.setAttribute('data-toast-wrapper', '');
        document.body.append(toastWrapper);
    }
    return toastWrapper;
}

function _buildToastBody(text) {
    var toastBodyWrapper = document.createElement('div');
    toastBodyWrapper.classList.add('d-flex'); 

    var toastBody = document.createElement('div');
    toastBody.setAttribute('class', 'toast-body');

    var img = document.createElement('img');
    img.setAttribute('class', 'rounded me-2');
    img.setAttribute('src', '');
    img.setAttribute('alt', '');

    var closeButton = document.createElement('button');
    closeButton.setAttribute('type', 'button');
    closeButton.setAttribute('class', 'btn-close btn-close-white me-2 m-auto');
    closeButton.setAttribute('data-bs-dismiss', 'toast');
    closeButton.setAttribute('data-label', 'Close');
   
    toastBody.textContent = text;
    toastBodyWrapper.append(toastBody);
    toastBodyWrapper.append(closeButton);

    return toastBodyWrapper;
}

// bgColorClass 'bg-success', 'bg-primary', 'bg-danger'
function _buildToast(text, bgColorClass, delayTime) {
    var toast = document.createElement('div');

    toast.setAttribute('class', 'toast align-items-center text-white border-0');
    toast.classList.add(bgColorClass); 
    toast.setAttribute('role', 'alert');
    toast.setAttribute('aria-live', 'assertive');
    toast.setAttribute('aria-atomic', 'true');
    toast.setAttribute('data-bs-delay', delayTime);
    
    var toastBody = _buildToastBody(text);
    toast.append(toastBody);
    return toast;
}