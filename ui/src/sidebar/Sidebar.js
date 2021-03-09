import React, {Component} from 'react'
import UsersList from "./UsersList";
import ChannelsList from "./ChannelsList";
import ContextSwitcher from "./ContextSwitcher";
import './Sidebar.css'
import Profile from "./Profile";
import {Route, Switch} from "react-router-dom";
import {observer} from "mobx-react";
import {StateContext} from "../store/StateContext";


class Sidebar extends Component {
    state = {context: ''}

    componentDidMount() {

    }

    updateContext = (context) => {
        this.setState({context});
    }

    render() {
        let {context} = this.state;
        let {location} = this.props;
        let {activeChannelId} = this.context
        let path = (location && location.pathname) || '';

        let ctxMatch = path.match(/\/(direct|channel|profile)(?:\/([^/]+))?/);
        if (ctxMatch) {
            if (!activeChannelId && ctxMatch[2]) {
                if (ctxMatch[1] === 'direct') {
                    this.context.setActiveChannelId('@' + ctxMatch[2]);
                } else if (ctxMatch[1] === 'channel') {
                    this.context.setActiveChannelId(parseInt(ctxMatch[2]));
                }
            }

            // this.state.context =  ctxMatch[1]
            context  = ctxMatch ? ctxMatch[1] : 'direct';
        }
        return (
            <div className='usersSidebar'>
                <ContextSwitcher onChange={this.updateContext} context={context}/>
                <Switch>
                    <Route path="/direct" component={UsersList}/>
                    <Route path="/channel" component={ChannelsList}/>
                    <Route path="/profile" component={Profile}/>
                    <Route path="/profile/:nickname" component={Profile}/>
                </Switch>
            </div>
        )
    }
}

Sidebar.contextType = StateContext

export default observer(Sidebar)