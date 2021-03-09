import {makeAutoObservable} from "mobx";


export default class Message {
    cid = null;
    mid = null;
    pid = null;
    uid = null;
    user = null;
    body = null;
    isThread = false;
    at = null;
    images = [];
    likes = [];
    self =  false;

    update = props => {
        let {pid, uid, user, body, isThread, created, images, likes, self} = props
        if (pid !== undefined) this.pid = pid;
        if (uid !== undefined) this.uid = uid;
        if (user !== undefined) this.user = user;
        if (body !== undefined) this.body = body;
        if (isThread !== undefined) this.isThread = isThread;
        if (images !== undefined) this.images = images || this.images
        if (likes !== undefined) this.likes = likes || this.likes
        if (self !== undefined) this.self = self
        if (created !== undefined) {
            this.at = Date.parse(created);
        }
    }

    constructor(props) {
        this.cid = props.cid;
        this.mid = props.mid;
        this.update(props)

        makeAutoObservable(this)
    }

}