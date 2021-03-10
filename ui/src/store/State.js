import {action, makeObservable, observable, runInAction} from 'mobx'
import User from "./User";
import Channel from "./Channel";
import Message from "./Message";


class State {

    // @var Socket
    socket = null;
    pingTimer = null;
    pingEnabled = true;
    lastPing = null;
    isOnline = false;
    user = null;
    token = null;
    users = {};
    channels = [];
    directs = {};
    activeChannelId = null;
    activeChannel = null;
    channelFixed = null;
    loaded = false
    channelsLoaded = false
    usersLoaded = false
    chatLoading = false
    chatLoaded = false


    startPing = (period = 1000) => {
        if (this.pingTimer === null) {
            this.pingTimer = setInterval(() => {
                if (this.pingEnabled) {
                    this.socket.ping().then(ping => {
                        this.setLastPing(ping)
                    })
                }
            }, period)
        }
    }

    stopPing = () => {
        if (this.pingTimer !== null) {
            clearInterval(this.pingTimer)
            this.pingTimer = null
        }
    }

    setLastPing = (ping) => {
        this.lastPing = ping
    }

    setOnline = (isOnline = false) => {
        this.isOnline = isOnline
    }

    setUser = (userData = {}) => {
        this.user = userData ? new User(userData) : null;
    }

    setChannels = (channelsData = []) => {
        // this.channels = channelsData
        this.channels = channelsData.map(channelData => {
            return channelData instanceof Channel ? channelsData : new Channel(channelData);
        })
    }

    setUsers = (users = {}) => {
        this.users = users;
    }

    getUsers = async (filter = {}) => {
        if (this.usersLoaded) {
            return
        }
        let {mappedUsers} = await this.socket.usersGetList(filter)
        this.mergeUsers(mappedUsers)
        this.usersLoaded = true
    }

    getUser = (nickname) => {
        return Object.values(this.users).find(user => user.nickname === nickname);
    }

    getUserByUid = (uid) => {
        return Object.values(this.users).find(user => user.uid === uid);
    }

    isUserOnline = (uid) => {
        let user = this.getUserByUid(uid);
        return user ? user.isOnline() : false;
    }

    getUsersStatus = (uids = []) => {
        let statuses = {}
        uids.forEach(uid => {
            statuses[uid] = false
            if (this.users[uid] && this.users[uid].isOnline) {
                statuses[uid] = this.users[uid].isOnline()
            } else {
                console.log(this.users[uid])
            }
        })
        return statuses
    }

    setActiveChannelId = (cid) => {
        this.activeChannelId = cid;
        this.getActiveChannel(cid).then((chan) => {
            runInAction(() => {
                this.activeChannel = chan
            })
        })
    }

    getActiveChannel = async (cid) => {
        if (typeof cid !== 'number' && typeof cid !== 'string') {
            return null;
        }
        if (cid.toString().substr(0, 1) === '@') {
            return await this.getDirect(cid.toString().substr(1));
        } else {
            return await this.getChannel(cid);
        }
    }

    updateActiveChannel = (channelData = {}) => {
        if (this.activeChannel) {
            this.activeChannel.update(channelData);
        }
    }

    getChannel = cid => {
        let chan = this.channels.find(channel => channel.cid === cid);
        if (chan && !chan.loaded) {
            this.getChannelMessages(chan.cid)
            this.getChannelMembers(chan.cid)
            this.mergeChannels([chan])
        }
        return chan
    }

    getDirect = async nickname => {
        let cid = this.directs[nickname]
        if (cid) {
            return await this.getChannel(cid);
        }
        let data = await this.socket.channelGetDirect({nickname});
        let directChan = data.channel || null;

        if (!directChan) {
            return null
        }
        this.directs[nickname] = directChan.cid
        runInAction(() => {
            this.mergeChannels([directChan])
        })
        return this.getChannel(directChan.cid)
    }

    getChannels = async (filter = {}) => {
        if (this.channelsLoaded) {
            return
        }
        let {channels} = await this.socket.channelsGetList(filter)
        this.channels = channels
        this.channelsLoaded = true
        return channels
    }

    createChannel = ({title}) => {
        if (this.user) {
            this.socket.channelCreate({title}).then((channel) => {
                this.mergeChannels([new Channel(channel)])
                this.setActiveChannelId(channel.cid);
            })
        }
    }

    dropChannel = () => {
        if (this.activeChannel && this.user) {
            let cid = this.activeChannel.cid;
            this.socket.channelDelete({cid}).then(() => {
                this.channels = (this.channels || []).filter(channel => channel.cid !== cid)
                this.setActiveChannelId(null);
            })
        }
    }

    joinChannel = () => {
        if (this.activeChannel && this.user) {
            let chan = this.activeChannel
            let uid = this.user.uid;
            this.socket.channelJoin({cid: chan.cid, uid}).then(({chan}) => {
                this.mergeChannels([chan])
            })
        }
    }

    leaveChannel = () => {
        if (this.activeChannel && this.user) {
            let chan = this.activeChannel
            let uid = this.user.uid;
            this.socket.channelLeave({cid: chan.cid, uid}).then(({chan}) => {
                this.mergeChannels([chan])
            })
        }
    }

    getChannelMembers = (cid) => {
        this.socket.channelGetMembers(cid).then(({members}) => {
            this.updateChannel(cid, {users: members})
        })
    }

    getChannelMessages = (cid, filter = {}) => {
        filter['cid'] = cid
        return this.getMessages(filter)
    }

    getMessages = (filter = {}) => {
        this.socket.messagesGetList(filter).then(({messages, users}) => {
            this.mergeUsers(users);
            this.mergeMessages(messages);
            return {messages, users}
        })
    }

    /**
     * @param {Channel[]} channels
     */
    mergeChannels = (channels = []) => {
        channels.forEach(chan => {
            if (!chan || !chan.cid) {
                return
            }
            let cur = this.channels.find(cc => cc.cid === chan.cid)
            if (cur) {
                cur.update(chan)
            } else {
                this.channels.push(chan)
            }
        })
    }

    mergeUsers = (users = []) => {
        Object.values(users).forEach(user => {
            this.users[user.uid] = user
        })
    }

    mergeMessages = (messages = []) => {
        let chanMap = {}
        messages.forEach(msg => {
            msg.user = msg.uid ? this.users[msg.uid] : null
            let [cid, mid] = [msg.cid, msg.mid]

            if (!chanMap[cid]) {
                chanMap[cid] = {};
                let chan = Object.values(this.channels).find(ch => ch.cid === cid) || {};
                (chan.messages || []).forEach(msg => {
                    chanMap[cid][msg.mid] = msg;
                })
            }

            chanMap[cid][mid] = msg
        })

        Object.entries(chanMap).forEach(([cid = "", chanMessages = {}]) => {
            let messages = Object.values(chanMessages)
            messages.sort((a, b) => {
                if (a.at < b.at) {
                    return -1;
                } else if (a.at > b.at) {
                    return 1;
                } else {
                    return 0;
                }
            })
            this.updateChannel(cid, {messages})
        })
    }

    updateChannel = (cid, data) => {
        cid = parseInt(cid)
        let chan = this.channels.find(chan => {
            return chan.cid === cid
        })

        if (!chan) {
            this.channels.push(new Channel({cid, ...data}));
        } else {
            chan.update(data)
        }
    }

    addMessage = async (body) => {
        let chan = this.activeChannel
        let user = this.user
        if (chan && user && body) {
            let msg = await this.socket.messagesCreate({cid: chan.cid, uid: user.uid, user, body,})
            this.mergeMessages([new Message(msg)])
            return msg
        } else {
            return null
        }
    }

    likeMessage = (cid, mid) => {
        let chan = this.activeChannel
        if (this.user && chan && chan.cid === cid) {
            const uid = this.user.uid
            let msg = chan.messages.filter(m => mid === m.id && m.cid === cid)
            if (!msg) {
                return
            }
            let likes = msg.likes || []
            let alreadyLiked = likes.filter(like => like === uid).length > 0
            console.log(msg)
            console.log(alreadyLiked)
            console.log(likes)
            if (alreadyLiked) {
                likes = likes.filter(like => like !== uid)
            } else {
                likes.push(uid)
            }
            console.log(likes)
            msg.likes = likes
        }
    }

    signInGoogle = ({profileId, accessToken, tokenId}) => {
        return this.socket.signInGoogle({
            profileId: profileId,
            accessToken: accessToken,
            secretToken: tokenId
        }).then(({me, token}) => {
            this.setUser(me)
            this.mergeUsers([me])
            this.setToken(token)
            return me
        })
    }

    signInTyrion = () => {
        return this.signInToken({token: "Tyrion"})
    }

    signInToken = ({token}) => {
        return this.socket.signInToken({token}).then(({me, token}) => {
            this.setUser(me)
            this.mergeUsers([me])
            this.setToken(token)
            return me
        })
    }

    signInTyrionAsync = async () => {
        return await this.signInToken({token: "Tyrion"})
    }

    signInTokenAsync = async ({token}) => {
        let {me, token: newToken} = await this.socket.signInToken({token})

        this.setUser(me)
        this.mergeUsers([me])
        this.setToken(newToken)

        return me
    }

    signOut = () => {
        this.socket.signOut().then(() => {
            this.resetAuth()
        })
    }

    resetAuth = () => {
        this.setUser(null);
        this.setToken(null);
    }

    setToken = (token) => {
        if (token) {
            this.token = token;
            window.sessionStorage.setItem('token', token);
        } else {
            this.token = null;
            window.sessionStorage.removeItem('token');
        }
    }

    getToken = () => {
        if (!this.token) {
            this.token = window.sessionStorage.getItem('token');
        }
        return this.token;
    }

    constructor({activeChannelId = null} = {}) {
        makeObservable(this, {
            pingEnabled: observable,
            lastPing: observable,
            isOnline: observable,
            user: observable,
            token: observable,
            users: observable,
            channels: observable,
            activeChannelId: observable,
            activeChannel: observable,
            channelFixed: observable,
            chatLoading: observable,
            chatLoaded: observable,
            getUser: observable,
            startPing: action,
            stopPing: action,
            setLastPing: action,
            setOnline: action,
            signInGoogle: action,
            signInToken: action,
            signInTyrion: action,
            signOut: action,
            resetAuth: action,
            getToken: action,
            setToken: action,
            setUser: action,
            getUsers: action,
            setUsers: action,
            mergeUsers: action,
            isUserOnline: action,
            getChannels: action,
            setChannels: action,
            mergeChannels: action,
            setActiveChannelId: action,
            updateActiveChannel: action,
            createChannel: action,
            dropChannel: action,
            joinChannel: action,
            updateChannel: action,
            leaveChannel: action,
            addMessage: action,
            likeMessage: action,
            mergeMessages: action,
        });

        if (!isNaN(activeChannelId)) {
            this.activeChannelId = activeChannelId;
            this.channelFixed = true;
        }
    }
}

export default State;
