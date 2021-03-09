import {makeAutoObservable} from 'mobx'
import Message from "./Message";

export default class Channel {
    cid = null
    title = "";
    owner = null;
    lastMessage = null;
    type = "";
    users = [];
    messages = [];

    update = ({title, users, message, messages, lastMessage, type}) => {
        if (title) this.title = title;
        if (users) this.users = users;//.map(userData => new User(userData));
        if (message) this.messages.push(new Message(message));
        if (lastMessage) this.lastMessage = lastMessage
        if (type) this.type = type
        if (messages) {
            // const currentLastMessage = this.messages[this.messages.length - 1] || {id: -1};
            // this.messages = this.messages.concat(
            //     messages.filter(message => message.id > currentLastMessage.id).map(message => new Message(message))
            // );
            this.messages = messages
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
        }  else if (this.isDirect()) {
            if (this.users.length === 1) return true; // self channel
            console.log(this.users)
            return (this.users||[]).filter(u => u !== user.uid &&  usersOnline[u]).length > 0
        }
        return false
    }

    constructor({cid, title, users, messages, owner, type}) {
        this.cid = cid;
        this.owner = owner;

        this.update({title, users, messages, type});
        makeAutoObservable(this)
    }

}


