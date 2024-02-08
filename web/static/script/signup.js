function Toggle() {
    let signin = document.getElementById('login-left-content')
    let signup = document.getElementById('Signup-left-content')
    if (signin.style.display === 'none') {
        signin.style.display = 'flex'
        signup.style.display = 'none'
    } else {
        signup.style.display = 'flex'
        signin.style.display = 'none'
    }
}


const modalContainer = document.querySelector(".modal-container");
const modalTrigger = document.querySelectorAll(".modal-trigger");
modalTrigger.forEach(trigger => trigger.addEventListener("click",toggleModal));
function toggleModal(){
    modalContainer.classList.toggle("active");
}
function toggleCategorie(){
    const parent = document.querySelector('.Category-checkbox');
    //on verifie si l'element est deja visible
    const isVisible = window.getComputedStyle(parent).display !== 'none'; 
    parent.style.display = isVisible ? 'none' : 'flex';
}