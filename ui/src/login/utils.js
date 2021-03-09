
import {createBrowserHistory} from "history";

let history = createBrowserHistory()

function GotoGoogleSignIn() {
    return (e) => history.replace(GoogleSignInUrl())
}

function GotoRoot() {
    return (e) => history.replace('/ui')
}


function GoogleSignInUrl() {
    return '/api/identity/login/google/?referrer=login'
}

function AccountSignInUrl() {
    return '/api/identity/login/account/?referrer=login'
}

function AccountSignUpUrl() {
    return '/api/identity/login/signup/?referrer=login'
}

export {
    GotoRoot,
    GotoGoogleSignIn,
    GoogleSignInUrl,
    AccountSignInUrl,
    AccountSignUpUrl,
}