function showSignup() {
    document.getElementById('login-left-content').style.display = 'none';
    document.getElementById('Signup-left-content').style.display = 'flex';
}

function showSignin() {
    document.getElementById('Signup-left-content').style.display = 'none';
    document.getElementById('login-left-content').style.display = 'flex';
}


const modalContainer = document.querySelector(".modal-container");
const modalTrigger = document.querySelectorAll(".modal-trigger");
const commentContainer = document.querySelector(".main-comment");
const commentTrigger = document.querySelectorAll(".comment-trigger");

commentTrigger.forEach(trigger => trigger.addEventListener("click",commentModal));
modalTrigger.forEach(trigger => trigger.addEventListener("click",toggleModal));
function commentModal(){
    commentContainer.classList.toggle("active1");
}
function toggleModal(){
    modalContainer.classList.toggle("active");
}

// Fonction pour stocker la position du défilement de la div centerside
function saveScrollPosition() {
    var centerside = document.getElementById('centerside');
    localStorage.setItem('centersideScrollPosition', centerside.scrollTop);
}

// Fonction pour charger et appliquer la position de défilement de la div centerside
function restoreScrollPosition() {
    var centerside = document.getElementById('centerside');
    var scrollPosition = localStorage.getItem('centersideScrollPosition');
    if (scrollPosition) {
        centerside.scrollTop = scrollPosition;
    }
}

// Écouteur d'événement pour enregistrer la position de défilement lors du défilement
document.getElementById('centerside').addEventListener('scroll', saveScrollPosition);

// Appel de la fonction pour restaurer la position de défilement lors du chargement de la page
window.addEventListener('load', restoreScrollPosition);

document.getElementById('file-input').addEventListener('change', function() {
    document.querySelector('.Photo').style.border = '2px solid #519e7a';
});