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