import Channel from "./Channel";
import User from "./User";
import Message from "./Message";
import Ping from "./Ping";
import Config from "./Config";

class Socket {
    eventPrefix = ""
    eventSec = 0
    state = null
    ws = null
    io = null // stub
    lastTry = null
    tries = 0
    pingSeq = 0
    onceQueue = {}
    listeners = {}

    constructor({state, server}) {
        this.server = server || window.location.origin.replace(/^http(s?):\/\//, 'ws$1://') + '/ws'
        this.state = state;
        this.eventPrefix = 'ws-client-' + Math.round(Math.random() * 10000).toString()
        this.on('auth.required', this.onAuthRequired)
        this.on('auth.signed-out', this.onSignedOut)
        this.on('ws.connected', this.noop)
        this.on('ws.disconnected', this.noop)
        this.on('users.info', this.onUserInfo)
        this.on('messages.created', this.onMessageCreated)
        this.on('channels.members', this.onChannelMembers)
        this.getWs().then(async () => {
            await this.state.loadConfig()
        })
    }

    getWs = () => {
        // let reconnectTimer = null
        return new Promise((resolve, reject) => {
            let reconnect = () => {
                if (this.ws && this.lastTry && Date.now() - this.lastTry < 5 * 1000 * 1000) {
                    return
                }
                let ws
                try {
                    this.tries++
                    this.lastTry = Date.now()
                    ws = new WebSocket(this.server);
                } catch (e) {
                    console.error(e)
                    return;
                }

                ws.onopen = (e) => {
                    this.lastTry = null
                    this.tries = 0
                    console.log('ws connected ', ws)
                    resolve(ws)
                    this.state.setOnline(true)
                    this.state.signIn()
                };

                ws.onmessage = e => {
                    // console.log(e)
                    this.dispatchResponse(e.data)
                }

                ws.onerror = (err) => {
                    this.state.setOnline(false)
                    this.ws = null
                    console.log('ws err', err)
                    ws.close()
                }

                ws.onclose = () => {
                    this.ws = null
                    console.log('ws disconnected')
                    this.state.setOnline(false)
                    setTimeout(() => {
                        reconnect();
                    }, 500);
                }

                this.ws = ws
            }

            let waitTime = Date.now() + 2000
            while (this.ws && this.ws.readyState === WebSocket.CONNECTING) {
                if (Date.now() > waitTime) {
                    break
                }
            }

            if (this.ws === null || this.ws.readyState === WebSocket.CLOSED || this.ws.readyState === WebSocket.CLOSING) {
                console.log('ws reconnect', this.ws)
                reconnect()
            } else if (this.ws.readyState === WebSocket.OPEN) {
                resolve(this.ws)
            }

            // setInterval(() => {
            //     if (this.ws === null || this.ws.readyState === WebSocket.CLOSED || this.ws.readyState === WebSocket.CLOSING) {
            //         reconnect()
            //     }
            // }, 5000)
        })
    }

    dispatchResponse = (msgs) => {
        msgs.split("\n").forEach(msg => {
            if (msg) {
                this.dispatchMsg(msg)
            }
        })
    }

    dispatchMsg = (msg) => {
        let raw = JSON.parse(msg)
        let {type, prev} = raw
        if (this.onceDequeue(type, raw, prev)) {
            return;
        } else if (this.notify(type, raw, prev)) {
            return;
        }
        console.log(msg)
    }

    send = async (type, payload) => {
        return this.getWs().then(ws => {
            let event = {type, payload, id: this.eventPrefix + (this.eventSec++)}
            console.log(['send', event.type || '???' , event, ws.readyState,])
            ws.send(JSON.stringify(event))
            return {event}
        })
    }

    on = (type, callback) => {
        if (!this.listeners[type]) {
            this.listeners[type] = [];
        }
        this.listeners[type].push(callback)
    }

    off = (type, callback) => {
        if (this.listeners[type]) {
            this.listeners[type] = this.listeners[type].filter(cb => cb !== callback)
        }
    }

    notify = (type, payload, prevId) => {
        let notified = false;
        (this.listeners[type] || []).forEach(cb => {
            cb(type, payload, prevId)
            console.log(['recv.on', type, prevId, payload])
            notified = true
        })
        return notified
    }

    once = (type, listener, prev = null) => {
        if (!this.onceQueue[type]) {
            this.onceQueue[type] = []
        }
        this.onceQueue[type].push([listener, prev])
    }

    onceDequeue = (type, payload, prevId) => {
        if (this.onceQueue[type] && this.onceQueue[type].length) {
            if (prevId) {
                this.onceQueue[type].filter(([, prev]) => prevId === prev).forEach(([listener,]) => {
                    console.log(['recv.once', type, prevId, payload])
                    listener(payload)
                })
                this.onceQueue[type] = this.onceQueue[type].filter(([, prev]) => prevId !== prev)
            }
            this.onceQueue[type].filter(([, prev]) => prev === null).forEach(([listener,]) => {
                console.log(['recv.once', type, null, payload])
                listener(payload)
            })
            this.onceQueue[type] = this.onceQueue[type].filter(([, prev]) => prev !== null)

            return true
        }
        return false
    }

    noop = () => {
        // do nothing
    }

    ping = async () => {
        let now = Date.now()
        let {event} = await this.send('ping')
        return new Promise((resolve, reject) => {
            this.once('pong', () => {
                resolve(new Ping(this.pingSeq++, now, Date.now()))
            }, event.id)
        })
    }

    getConfig = async () => {
        let {event} = await this.send('system.get-config', null)
        return new Promise((resolve, reject) => {
            this.once('system.config', (data) => {
                let payload = data.payload || {};
                if (!payload.ok) {
                    reject(payload.message || 'error')
                    return
                }
                resolve(new Config(payload.config || {}));
            }, event.id)
        })
    }

    signInGoogle = async ({accessToken, secretToken, ttl = 3600}) => {
        let payload = {service: 'google', accessToken, secretToken, ttl}
        let {event} = await this.send('auth.sign-in', payload)
        return this.onSignedIn(event)
    }

    signInToken = async ({token}) => {
        let payload = {service: 'token', accessToken: token}
        let {event} = await this.send('auth.sign-in', payload)
        return this.onSignedIn(event)
    }

    onSignedIn = event => {
        return new Promise((resolve, reject) => {
            this.once('auth.signed-in', (data) => {
                let payload = data.payload || {};
                if (!payload.ok) {
                    reject(payload.message || 'error')
                    return
                }
                let {me, token} = payload
                if (me) {
                    resolve({me, token});
                } else {
                    reject("auth failed")
                }
            }, event.id)
        })
    }

    signOut = async () => {
        let {event} = await this.send('auth.sign-out', null)
        return event
    }

    onAuthRequired = (type, event, prevId) => {
        this.state.resetAuth()
    }

    onSignedOut = () => {
        this.state.resetAuth()
    }

    onUserInfo = (type, event, prevId) => {
        if (type === 'users.info') {
            let user = event.payload && event.payload.user
            if (user) {
                this.state.mergeUsers([new User(user)])
            }
        }
    }

    onChannelMembers = (type, event, prevId) => {
        if (type === 'channels.members') {
            let payload = event.payload || {};
            let cid = payload.cid || null;
            let members = payload.members || [];
            if (cid) {
                this.state.updateChannel(cid, {users: members})
            }
        }
    }
    onMessageCreated = (type, event, prevId) => {
        if (type === 'messages.created') {
            let payload = event.payload || {};
            let msg = payload.msg || null;
            if (msg) {
                this.state.mergeMessages([new Message(msg)])
            }
        }
    }

    channelsGetList = async (filter) => {
        let {event} = await this.send('channels.get-list', filter ? filter : {})
        return new Promise((resolve, reject) => {
            this.once('channels.list', (data) => {
                let payload = data.payload || {};
                if (!payload.ok) {
                    reject(payload.message || 'error')
                    return
                }
                let channels = payload.channels || []
                channels = channels.map(channelData => new Channel(channelData));
                resolve({channels});
            }, event.id)
        })
    }

    channelGetDirect = async ({uid, nickname}) => {
        let {event} = await this.send('channels.get-direct', {nickname, uid})
        return this.singleChanPromise('channels.direct', event)
    }

    channelCreate = async ({title}) => {
        let {event} = await this.send('channels.create', {title})
        return this.singleChanPromise('channels.created', event)
    }

    channelDelete = async ({cid}) => {
        let {event} = await this.send('channels.delete', {cid})
        return this.singleChanPromise('channels.deleted', event)
    }

    channelJoin = async ({cid, uid}) => {
        let {event} = await this.send('channels.join', {cid, uid})
        return this.singleChanPromise('channels.joined', event)
    }

    channelLeave = async ({cid, uid}) => {
        let {event} = await this.send('channels.leave', {cid, uid})
        return this.singleChanPromise('channels.left', event)
    }

    channelSetLastSeen = async ({cid, mid}) => {
        let {event} = await this.send('channels.set-last-seen', {cid, mid})
        return event.payload ? event.payload.ok : event.payload.message
    }

    channelGetLastSeen = async ({cid}) => {
        let {event} = await this.send('channels.get-last-seen', {cid})
        return new Promise((resolve, reject) => {
            this.once('channels.last-seen', (data) => {
                let payload = data.payload || {};
                if (!payload.ok) {
                    reject(payload.message || 'error');
                    return;
                }
                let {mid} = data.payload;
                resolve({mid});
            }, event.id)
        })
    }

    singleChanPromise = (eventType, prevEvent) => {
        return new Promise((resolve, reject) => {
            this.once(eventType, (data) => {
                let payload = data.payload || {};
                if (!payload.ok) {
                    reject(payload.message || 'error')
                    return
                }
                let channelData = payload.channel || null
                if (channelData) {
                    resolve({channel: new Channel(channelData)});
                }
            }, prevEvent.id)
        })
    }

    channelGetMembers = async (cid) => {
        let {event} = await this.send('channels.get-members', {cid})
        return new Promise((resolve, reject) => {
            this.once('channels.members', (data) => {
                let payload = data.payload || {};
                if (!payload.ok) {
                    reject(payload.message || 'error');
                    return;
                }
                let {members} = data.payload;
                resolve({members});
            }, event.id)
        })
    }

    usersGetList = async (filter) => {
        let {event} = await this.send('users.get-list', filter ? filter : {})

        return new Promise((resolve, reject) => {
            this.once('users.list', (data) => {
                let payload = data.payload || {};
                if (!payload.ok) {
                    reject(payload.message || 'error')
                    return
                }
                let users = payload.users || [];
                let mappedUsers = {};
                users.forEach(userData => {
                    mappedUsers[userData.uid] = new User(userData)
                });
                resolve({mappedUsers})
            }, event.id)
        })
    }

    messagesGetList = async (filter) => {
        let {event} = await this.send('messages.get-list', filter || {})
        return new Promise((resolve, reject) => {
            this.once('messages.list', (data) => {
                let payload = data.payload || {};
                if (!payload.ok) {
                    reject(payload.message || 'error');
                    return;
                }

                let {messages, users} = data.payload;
                users = Object.values(users || {}).map(userData => new User(userData));
                messages = (messages || []).map(msgData => new Message(msgData));
                resolve({messages, users});
            }, event.id)
        })
    }

    messagesCreate = async ({cid, uid, pid, body,}) => {
        let {event} = await this.send('messages.create', {cid, uid, pid, body})
        return new Promise((resolve, reject) => {
            this.once('messages.created', (data) => {
                let payload = data.payload || {};
                if (!payload.ok) {
                    reject(payload.message || 'error');
                    return;
                }
                let msg = payload.msg || null;
                if (msg) {
                    resolve(msg);
                }
            }, event.id)
        })
    }


}

export default Socket;
