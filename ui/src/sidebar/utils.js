
import {createBrowserHistory} from "history";

let history = createBrowserHistory()

function GotoChannel(chan) {
    return (e) => history.push(ChannelUrl(chan))
}

function ChannelUrl(chan) {
    return chan ? '/ui/channel/'+chan.cid : '/ui/channel'
}

function ChannelLink(chan) {
    return chan ? '/channel/'+chan.cid : '/ui/channel'
}

function GotoDirect(user) {
    return (e) => history.push(DirectUrl(user))
}


function DirectUrl(user) {
    return user ? '/ui/direct/'+(user.nickname || user.uid) : '/ui/direct'
}

function ProfileLink(user) {
    return user ? '/profile/'+(user.nickname || user.uid) : '/ui/direct'
}

function GotoLogout(user) {
    return (e) => history.push('/ui/logout')
}

function GotoProfile(user) {
    return (e) => history.push('/ui/profile')
}

export {
    GotoChannel,
    ChannelUrl,
    ChannelLink,
    GotoDirect,
    DirectUrl,
    GotoLogout,
    GotoProfile,
    ProfileLink,
}