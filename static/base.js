
// Login input

document.addEventListener("DOMContentLoaded", function () {
    const loginForm = document.getElementById("loginForm")

    if (loginForm) {
        LoginJS()
    }


});

function LoginJS(){
    const emailInput = document.getElementById("email");
    const passwordInput = document.getElementById("password");
    const loginButton = document.getElementById("login");
    const form = document.getElementById("loginForm");

    // regular expression for email format
    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

    function validateForm() {
        const isEmailValid = emailPattern.test(emailInput.value);
        const isPasswordValid = passwordInput.value.trim().length > 0;

        (!isEmailValid && !isPasswordValid) 
        ? loginButton.classList.add("btn-disabled") 
        : loginButton.classList.remove("btn-disabled")
    }

    // Validate in real time inputs
    emailInput.addEventListener("input", validateForm);
    passwordInput.addEventListener("input", validateForm);

    
    form.addEventListener("submit", function (event) {
        event.preventDefault()

        // if inputs values are not correct, send alert
        if (!emailPattern.test(emailInput.value)) {
            alert("Por favor, introduce un correo válido.")
        }

        if (passwordInput.value.trim().length == 0){
            alert("Por favor, introduce una contraseña.")
        }
        
        fetch("/api/login", {
            method: "POST",
            body: JSON.stringify({
                email: emailInput.value,
                password: passwordInput.value
            })
        })
        .then(res => {
            if(res.status == 404){
                DisplayAlert("No existe un usuario registrado con el correo proporcionado")
            }else if(res.status == 400){
                DisplayAlert("La contraseña proporcionada no es correcta")
            }else if(res.status == 500) {
                DisplayAlert("Hubo un error en los servidores. Por favor, intenta más tarde")
            }else if(res.status == 200){
                // Redirigir a home
                window.location.replace("/home")
            }
        })

    })
}

function DisplayAlert(msg) {
    let alertContainer = document.getElementById("alert")
    const alertHTML = `
    <div class="alert alert-danger d-flex align-items-center" role="alert">
        <svg class="bi flex-shrink-0 me-2" width="24" height="24" role="img" aria-label="Danger:"><use xlink:href="#exclamation-triangle-fill"/></svg>
        <div>
            ${msg}
        </div>
        <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
    </div>`;

    alertContainer.innerHTML = alertHTML
    setTimeout(() => {
        DisableAlert()
    }, 5000)
}

function DisableAlert(){
    let alertContainer = document.getElementById("alert")
    alertContainer.innerHTML = ""
}