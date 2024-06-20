
let switchCtn = document.querySelector("#switch-cnt");
let switchC1 = document.querySelector("#switch-c1");
let switchC2 = document.querySelector("#switch-c2");
let switchCircle = document.querySelectorAll(".switch__circle");
let switchBtn = document.querySelectorAll(".switch-btn");
let aContainer = document.querySelector("#a-container");
let bContainer = document.querySelector("#b-container");
let allButtons = document.querySelectorAll(".submit");

let getButtons = (e) => e.preventDefault()

let changeForm = (e) => {

    switchCtn.classList.add("is-gx");
    setTimeout(function(){
        switchCtn.classList.remove("is-gx");
    }, 1500)

    switchCtn.classList.toggle("is-txr");
    switchCircle[0].classList.toggle("is-txr");
    switchCircle[1].classList.toggle("is-txr");

    switchC1.classList.toggle("is-hidden");
    switchC2.classList.toggle("is-hidden");
    aContainer.classList.toggle("is-txl");
    bContainer.classList.toggle("is-txl");
    bContainer.classList.toggle("is-z200");
}

let mainF = (e) => {
    for (var i = 0; i < allButtons.length; i++)
        allButtons[i].addEventListener("click", getButtons );
    for (var i = 0; i < switchBtn.length; i++)
        switchBtn[i].addEventListener("click", changeForm)
}

window.addEventListener("load", mainF);

// Yorum popup elementlerini seçelim
const yorumPopup = document.getElementById("yorumPopup");
const yorumClose = document.querySelector(".yorum-popup-close");
const yorumEkleButton = document.getElementById("yorumEkleButton");
const yorumInput = document.getElementById("yorumInput");
const yorumListe = document.querySelector(".yorum-liste");

// "Yorum Yap" butonlarına event listener ekle
document.querySelectorAll('.post__footerButton .material-icons').forEach(button => {
    if (button.innerText === 'comment') {
        button.parentElement.addEventListener('click', () => {
            yorumPopup.style.display = "block";
        });
    }
});

// Popup kapama işlevi
yorumClose.onclick = function() {
    yorumPopup.style.display = "none";
}

// Popup dışına tıklanırsa kapat
window.onclick = function(event) {
    if (event.target == yorumPopup) {
        yorumPopup.style.display = "none";
    }
}

// Yorum ekleme butonu
yorumEkleButton.onclick = function() {
    const yorumMetni = yorumInput.value.trim();
    if (yorumMetni !== "") {
        const yeniYorum = document.createElement("p");
        yeniYorum.innerHTML = `<strong>Sen:</strong> ${yorumMetni}`;
        yorumListe.appendChild(yeniYorum);
        yorumInput.value = "";
    }
}
// Modal'ı açma fonksiyonu
function openModal() {
    document.getElementById("editProfileModal").style.display = "block";
}

// Modal'ı kapatma fonksiyonu
function closeModal() {
    document.getElementById("editProfileModal").style.display = "none";
}

// Kullanıcı modal dışında bir yere tıklarsa modal kapanır
window.onclick = function(event) {
    if (event.target == document.getElementById("editProfileModal")) {
        closeModal();
    }
}
