import axios from "axios";
import Channel from "./Channel";
import User from "./User";
import Message from "./Message";
import {object} from "prop-types";


export default function SeedGoT({state}) {
    state.setActiveChannelId(1)
    getUsers(state)
    getChannels(state)
    getMessages(state)
}

function getUsers(state) {
    axios.get('https://raw.githubusercontent.com/jeffreylancaster/game-of-thrones/master/data/characters.json').then(
        resp => resp.data
    ).then(
        data => {
            let chars = data.characters || []
            return chars.filter(c => c.characterImageFull).map((c, idx) => {
                return new User({
                    uid: idx,
                    nickname: nickname(c.characterName),
                    title: c.characterName,
                    avatarUrl: c.characterImageFull,
                    email: name2email(c.characterName).toLowerCase(),
                })
            })
        }
    ).then(users => {
        state.setUsers(users)
    })
}

function nickname(name) {
    return name.replace(/(\W+)/g, '')
}

function name2email(name) {
    return name.replace(/(\W+)/g, '.') + '@got.example';
}

function getChannels(state) {
    axios.get('https://raw.githubusercontent.com/jeffreylancaster/game-of-thrones/master/data/episodes.json').then(
        resp => resp.data
    ).then(
        data => {
            let episodes = data.episodes || []
            return episodes.map(e => {
                let cid = `S${e.seasonNum}E${e.episodeNum}`
                let title = `${cid}: ${e.episodeTitle}`
                return new Channel({cid, title})
            })
        }
    ).then(channels => {
        state.setChannels(channels)
    })
}

function getMessages(state) {
    axios.get('https://raw.githubusercontent.com/jeffreylancaster/game-of-thrones/master/data/script-bag-of-words.json').then(
        resp => resp.data
    ).then(
        data => {
            let episodes = data || []
            let usersMap = {}
            state.users.forEach(user => {
                usersMap[user.title] = user
            })
            let now = Math.floor(Date.now()/1000)
            let messages = {}, members = {};
            episodes.forEach(e => {
                let eMessages = [], eMembers = {};
                let items = e.text || []
                let first = now - items.length
                items.forEach((i,idx) => {
                    let user = usersMap[i.name]
                    if (user) {
                        eMembers[user.uid] = true
                        eMessages.push(new Message({
                            id: idx, user: user, content: i.text, self: user.nickname=== 'TyrionLannister', at: first+idx

                        }))
                    }
                })
                messages[e.episodeAlt] = eMessages
                members[e.episodeAlt] = Object.keys(eMembers)
            })
            return {messages, members}
        }
    ).then(({messages,members})=> {
        Object.values(state.channels).forEach(chan => {
            chan.messages = messages[chan.cid] || []
            chan.users = members[chan.cid] || []
        })
    })
}

function getChannelFeed(cid) {
    let messages = [
        {
            self: true,
            user: {
                title: "Jaime Lannister",
                nickname: "JaimeLannister",
                avatarUrl: 'https://images-na.ssl-images-amazon.com/images/M/MV5BMjIzMzU1NjM1MF5BMl5BanBnXkFtZTcwMzIxODg4OQ@@._V1_.jpg',
                uid: 1,
            },
            action: "added you as a friend",
            createdAt: "1 Hour Ago",
            text: `Ours is a life of constant reruns. We're always circling back to where
                    we'd we started, then starting all over again. Even if we don't run
                    extra laps that day, we surely will come back for more of the same
                    another day soon.`,
            images: [
                'https://react.semantic-ui.com/images/wireframe/image.png',
                'https://react.semantic-ui.com/images/wireframe/image.png',
            ],
            likes: [
                1, 2, 3, 4
            ],
        },
        {
            self: false,
            user: {
                title: "Helen Troy",
                nickname: "HelenTroy",
                avatarUrl: 'https://images-na.ssl-images-amazon.com/images/M/MV5BMjA2MDAwOTI0OV5BMl5BanBnXkFtZTcwNjA3NDg1Nw@@._V1_SX1500_CR0,0,1500,999_AL_.jpg',
                uid: 2,
            },
            action: " added 2 new illustrations",
            createdAt: "4 days ago",
            images: [
                'https://react.semantic-ui.com/images/wireframe/image.png',
                'https://react.semantic-ui.com/images/wireframe/image.png',
            ],
            likes: [
                1, 2, 3, 4, 5, 6
            ],
        },
        {
            self: false,
            user: {
                title: "Jenny Hess",
                nickname: "JennyHess",
                avatarUrl: 'https://images-na.ssl-images-amazon.com/images/M/MV5BMjA2MDAwOTI0OV5BMl5BanBnXkFtZTcwNjA3NDg1Nw@@._V1_SX1500_CR0,0,1500,999_AL_.jpg',
                uid: 3,
            },
            action: "added you as a friend",
            createdAt: "1 Hour Ago",
            likes: [
                1, 2, 3, 4
            ],
        },
        {
            self: true,
            user: {
                title: "Jaime Lannister",
                nickname: "JaimeLannister",
                avatarUrl: 'https://images-na.ssl-images-amazon.com/images/M/MV5BMjIzMzU1NjM1MF5BMl5BanBnXkFtZTcwMzIxODg4OQ@@._V1_.jpg',
                uid: 1,
            },
            action: "added you as a friend",
            createdAt: "1 Hour Ago",
            text: `Ours is a life of constant reruns. We're always circling back to where
                    we'd we started, then starting all over again. Even if we don't run
                    extra laps that day, we surely will come back for more of the same
                    another day soon.`,
            images: [
                'https://react.semantic-ui.com/images/wireframe/image.png',
                'https://react.semantic-ui.com/images/wireframe/image.png',
            ],
            likes: [
                1, 2, 3, 4
            ],
        },
        {
            self: false,
            user: {
                title: "Justen Kitsune",
                nickname: "JustenKitsune",
                avatarUrl: 'https://images-na.ssl-images-amazon.com/images/M/MV5BMjA2MDAwOTI0OV5BMl5BanBnXkFtZTcwNjA3NDg1Nw@@._V1_SX1500_CR0,0,1500,999_AL_.jpg',
                uid: 4,
            },
            action: "added you as a friend",
            createdAt: "1 Hour Ago",
            likes: [
                1, 2, 3, 4
            ],
        },
    ];
    return messages
}

