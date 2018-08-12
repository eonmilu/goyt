const TOKEN_AUTH_URL = "https://oxygenrain.com/yourtime/auth/token";
const DEFAULT_TIMEOUT = 1500;

function onSignIn(googleUser) {
    var idToken = googleUser.getAuthResponse().id_token;

    $.ajax({
        method: "POST",
        url: TOKEN_AUTH_URL,
        contentType: "application/x-www-form-urlencoded",
        data: {
            idtoken: idToken
        },
        timeout: DEFAULT_TIMEOUT
    }).done((response) => {
        console.log("Data succesfully posted");
        console.log(response);
    }).fail((jqXHR, textStatus, error) => {
        console.log(error);
        console.log(jqXHR);
        console.log(textStatus);
    });
}

function signOut() {
    var auth2 = gapi.auth2.getAuthInstance();
    auth2.signOut().then(function () {
        console.log('User signed out.');
    });
}
