import {makeAutoObservable} from 'mobx'
import Message from "./Message";

export default class Channel {
    cid = null
    title = "";
    owner = null;
    lastMessage = null;
    currentLastMessage = null;
    lastLoadedMessage = null;
    type = "";
    membersCount = 0;
    messagesCount = 0;
    users = [];
    messages = [];

    update = ({title, users, message, messages, lastMessage, type, membersCount, messagesCount, members_count, messages_count}) => {
        if (title) this.title = title;
        if (users) this.users = users;//.map(userData => new User(userData));
        if (message) this.messages.push(new Message(message));
        if (type) this.type = type
        if (membersCount) this.membersCount = membersCount
        if (messagesCount) this.messagesCount = messagesCount
        if (members_count) this.membersCount = members_count
        if (messages_count) this.messagesCount = messages_count
        if (type) this.type = type
        if (lastMessage) this.lastMessage = lastMessage
        if (messages) {
            const currentLastMessage = this.messages[this.messages.length - 1] || null;
            this.messages = currentLastMessage ? this.messages.concat(
                messages.filter(message => message.mid > currentLastMessage.id)
            ) : messages;
            this.messages = messages
            if (this.messages.length) {
                // this.lastLoadedMessage = this.messages[this.messages.length-1].mid
            }
        }
    }

    isMember = (user) => {
        let uid = user.uid
        return this.users ? this.users.filter(member => member === uid).length > 0 : false;
    }

    isPublic = () => this.type === 'public'
    isDirect = () => this.type === 'direct'

    isWritable = (user, usersOnline = {}) => {
        if (!this.cid || !user.uid) {
            return false
        } else if (this.isPublic()) {
            return this.isMember(user)
        } else if (this.isDirect()) {
            if (this.users.length === 1) return true; // self channel
            return (this.users || []).filter(u => u !== user.uid && usersOnline[u]).length > 0
        }
        return false
    }

    markSeen = (lastSeenId) => {
        this.messages.forEach(msg => {
            if (msg.mid <= lastSeenId) {
                msg.markSeen()
            }
        })
    }

    constructor(props) {
        let {cid,owner} = props
        this.cid = cid;
        this.owner = owner;

        this.update(props);
        makeAutoObservable(this)
    }

}


