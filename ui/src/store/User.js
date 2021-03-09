import {makeAutoObservable} from "mobx";

export default class User {
    uid = null;
    title = null;
    nickname = null;
    avatarUrl = null;
    email = null;
    status = "";

    update = props => {
        let {title, nickname, avatarUrl, email, status} = props
        this.title = title || nickname;
        this.nickname = nickname || title;
        this.avatarUrl = avatarUrl;
        this.email = email;
        this.status = status
    }

    isOnline = () => {
        return this.status === 'online'
    }

    constructor(props) {
        this.uid = props.uid || null;
        this.update(props)
        makeAutoObservable(this)
    }
}

