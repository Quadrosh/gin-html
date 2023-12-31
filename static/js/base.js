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


// Toast snackbar 
// parent class, no init in constructor
class Toast {
    constructor(text){
        // call it from the heirs
        // const toastElement = this.buildToast(text, 'bg-dark', 10000); // 'bg-success', 'bg-primary', 'bg-danger'
        // this.init(toastElement)
    }
    init(_toastElement){
        const toastWrapper = this.getOrCreateToastWrapper();
        toastWrapper.append(_toastElement);
        this.bootstrapToast = bootstrap.Toast.getOrCreateInstance(_toastElement);
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

    getOrCreateToastWrapper() {
        let toastWrapper = document.querySelector('body > [data-toast-wrapper]');
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
    
    buildToastBody(text) {
        const toastBodyWrapper = document.createElement('div');
        toastBodyWrapper.classList.add('d-flex'); 
    
        const toastBody = document.createElement('div');
        toastBody.setAttribute('class', 'toast-body');
    
        const img = document.createElement('img');
        img.setAttribute('class', 'rounded me-2');
        img.setAttribute('src', '');
        img.setAttribute('alt', '');
    
        const closeButton = document.createElement('button');
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
    buildToast(text, bgColorClass, delayTime) {
        const toast = document.createElement('div');
    
        toast.setAttribute('class', 'toast align-items-center text-white border-0');
        toast.classList.add(bgColorClass); 
        toast.setAttribute('role', 'alert');
        toast.setAttribute('aria-live', 'assertive');
        toast.setAttribute('aria-atomic', 'true');
        toast.setAttribute('data-bs-delay', delayTime);
        
        const toastBody = this.buildToastBody(text);
        toast.append(toastBody);
        return toast;
    }
   
}

/**
 * Example call 
 * (new ToastError('text')).show()
 *  */ 
class ToastError extends Toast {
    constructor(text){
        super()
        const toastElement = this.buildToast(text, 'bg-danger', 10000); // 'bg-success', 'bg-primary', 'bg-danger'
        this.init(toastElement)
    }
}


/**
 * Example call 
 * (new ToastInfo('text')).show()
 *  */ 

// ToastInfo snackbar 
class ToastInfo extends Toast{
    constructor(text){
        super()
        const toastElement = this.buildToast(text, 'bg-primary', 6000); 
        this.init(toastElement)
    }
}




/**
 * Example modall call 
 
(new Modal(
    'Осторожно',          // title
    'Дорогу переходят утки', // text
    'Закрыть',            // noBtnName
    'Да',                 // yesBtnName
    ()=>{alert('fuck')},  // yesBtnFunction
    )).show()

 *  */ 

// Modal constructor
class Modal {
    constructor(
        title, 
        text, 
        noBtnName='Закрыть', 
        yesBtnName='', 
        yesBtnAction=()=>{}
        ){
        const wrap = this.buildModal(title, text, noBtnName, yesBtnName, yesBtnAction)
        document.body.append(wrap)
        this.bsModal = bootstrap.Modal.getOrCreateInstance(wrap);
        this.show = function() {
            this.bsModal.show();
        }
    }


     buildModal(title, text, noBtnName, yesBtnName, yesBtnFunc) {
        const modal = document.createElement('div')
        modal.setAttribute('class', 'modal fade')
        modal.setAttribute('tabindex', '-1')
        modal.setAttribute('aria-labelledby', 'modalLabel')
        modal.setAttribute('aria-hidden', 'true')
        const modDialog = document.createElement('div')
        modDialog.setAttribute('class', 'modal-dialog')
        const modContent = document.createElement('div')
        modContent.setAttribute('class', 'modal-content')
        const header = this.buildModalHeader(title)
        modContent.append(header)
        const body = document.createElement('div')
        body.setAttribute('class', 'modal-body')
        body.innerText = text
        modContent.append(body)
        const footer = this.buildModalFooter(noBtnName, yesBtnName, yesBtnFunc)
        modContent.append(footer)
        modDialog.append(modContent)
        modal.append(modDialog)
        return modal
    }
    
    buildModalHeader(text) {
        const header = document.createElement('div');
        header.setAttribute('class', 'modal-header');
        header.setAttribute('style', 'border-bottom: none;');
    
        const title = document.createElement('h5');
        title.setAttribute('class', 'modal-title');
        title.setAttribute('id', 'modalLabel');
        title.innerText = text
    
        const closeBtn = document.createElement('button');
        closeBtn.setAttribute('class', 'btn-close');
        closeBtn.setAttribute('data-bs-dismiss', 'modal');
        closeBtn.setAttribute('aria-label', 'Close');
    
        header.append(title)
        header.append(closeBtn)
        return header
    }
    
     buildModalFooter(noBtnName, yesBtnName, yesBtnFunc) {
        const footer = document.createElement('div')
        footer.setAttribute('class', 'modal-footer')
        footer.setAttribute('style', 'border-top: none;')
    
        if (noBtnName){
            const noBtn = document.createElement('button')
            noBtn.setAttribute('type', 'button')
            noBtn.setAttribute('class', 'btn btn-secondary')
            noBtn.setAttribute('data-bs-dismiss', 'modal')
            noBtn.innerText = noBtnName
            footer.append(noBtn)
        } 
    
        if (yesBtnName && yesBtnFunc){
            const yesBtn = document.createElement('button')
            yesBtn.setAttribute('type', 'button')
            yesBtn.setAttribute('class', 'btn btn-primary')
            yesBtn.innerText = yesBtnName
            yesBtn.addEventListener('click', yesBtnFunc)
            footer.append(yesBtn)
        }
        return footer
    }
    
}
