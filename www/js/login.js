window.addEventListener("load", function () {
    const loginPath = "/login"

    function postLoginForm() {
        var xhr = new XMLHttpRequest();

        xhr.addEventListener("load", function(event) {
            console.log(`Response text ${event.target.responseText}`);
            if (event.target.status == 200)
                window.location = "/home";
        });

        xhr.addEventListener("error", function(event) {
            console.log(`Server return error message ${event.message}`);
        });

        xhr.open("POST", loginPath);
        xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");

        var formData = new FormData(form);
        console.log(`Post form to ${loginPath}`)
        xhr.send(`username=${formData.get("username")}&password=${formData.get("password")}&csrfToken=${formData.get("csrfToken")}`);
    }

    var form = document.getElementById("loginForm");
    form.addEventListener("submit", function (event) {
        event.preventDefault();
        postLoginForm();
    });
});