import React from "react";
import {Route} from "react-router";
import {observer} from "mobx-react";
import {BrowserRouter as Router} from "react-router-dom";
import {StateContext} from '../store/StateContext';
import ChatScreen from "../screen/ChatScreen";
import LoginScreen from "../screen/LoginScreen";
import './App.css';
import Dashboard from "../sys/dashboard";

class App extends React.Component {
    render() {
        let {user} = this.context;

        return (
            <div className="App">
                <Dashboard/>
                {user ? <AppRouter/> : <LoginScreen/>}
            </div>
        );
    }
}


function AppRouter() {
    return (
        <Router basename='ui'>
            <Route path="/" exact={true} component={ChatScreen}/>
            <Route path="/direct" component={ChatScreen}/>
            <Route path="/profile" component={ChatScreen}/>
            <Route path="/channel" component={ChatScreen}/>
        </Router>
    )
}

App.contextType = StateContext;

export default App = observer(App);
