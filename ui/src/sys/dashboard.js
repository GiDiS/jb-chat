import React from 'react'
import {StateContext} from "../store/StateContext";
import User from "../store/User";
import './sys.css'
import {observer} from "mobx-react";

class Dashboard extends React.Component {
    render() {
        let {user, isOnline, lastPing} = this.context;
        return (
            <div className='sysDashboard'>
                <ul>
                    <li><UserStatus user={user}/></li>
                    <li><OnlineStatus isOnline={isOnline}/></li>
                    <li><PingStatus ping={lastPing}/></li>
                </ul>
            </div>
        )
    }
}

Dashboard.contextType = StateContext

function UserStatus({user}) {
    return user instanceof User ? user.nickname || '?' : 'anon'
}

function OnlineStatus({isOnline}) {
    return isOnline ? "online" : "offline"
}

function PingStatus({ping}) {
    if (!ping) return null
    if (ping.finished) {
        let diff = (+ping.finished) - (+ping.started)
        return 'ping: ' + diff + 'ms'
    } else {
        return 'ping: failed'
    }
}


export default observer(Dashboard)