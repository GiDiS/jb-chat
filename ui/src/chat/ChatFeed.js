import React from 'react'
import {Button, Feed} from 'semantic-ui-react'
import './ChatFeed.css'
import {StateContext} from '../store/StateContext'
import {ProfileLink} from "../sidebar/utils";
import {observer} from "mobx-react";
import {Link} from "react-router-dom";

const ChatFeedEvent = (msg, msgIdx, self) => {
    let {images,user} = msg;
    // let likesCnt = likes.length
    let created = msg.at ? new Date(msg.at) : null;
    return user ? (
        <Feed.Event key={'msg-' + msgIdx} className={self ? 'self' : null}>
            <Feed.Label>
                <Link to={ProfileLink(user)}><img src={user.avatarUrl} alt={msg.user.title}/></Link>
            </Feed.Label>
            <Feed.Content>
                <Feed.Summary>
                    <Feed.User as='span'>
                        <Link to={ProfileLink(user)}>{msg.user.title}</Link>
                    </Feed.User>
                    {msg.action}

                </Feed.Summary>
                <Feed.Extra text>{msg.body}</Feed.Extra>
                {images ? (
                    <Feed.Extra images>{
                        msg.images.map((img, imgIdx) => (
                            <a href={img} key={'msg-img-' + imgIdx}><img src={img} alt={'img-' + imgIdx}/></a>))
                    }
                    </Feed.Extra>) : null
                }
                <Feed.Meta>
                    {created ? (<Feed.Date>{created.toLocaleString()}</Feed.Date>) : null}
                    {/*onLikeClick ? (<Feed.Like onClick={() => onLikeClick(msg)}> <Icon name='like'/>{likesCnt} Likes </Feed.Like>) : null*/}
                </Feed.Meta>
            </Feed.Content>
        </Feed.Event>
    ) : null;
}

class ChatFeed extends React.Component {
    state = {showScrollDown: false}

    constructor(props) {
        super(props)
        this.feedRef = React.createRef()
    }

    componentDidMount() {
        this.updateScrollDown()
    }

    updateScrollDown = () => {
        if (this.feedRef.current) {
            Object.values(this.feedRef.current.children).forEach(e => {
                if (e.classList.contains('feed')) {
                    if (e.clientHeight < e.scrollHeight && e.scrollTop < e.scrollHeight ) {
                        this.setState({showScrollDown: true})
                    }
                }
            })
        }
    }

    scrollDown = (e) => {
        Object.values(this.feedRef.current.children).forEach(e => {
            if (e.classList.contains('feed')) {
                if (e.clientHeight < e.scrollHeight && e.scrollTop < e.scrollHeight ) {
                    e.scrollTop = e.scrollHeight
                    this.setState({showScrollDown: false})
                }
            }
        })
    }

    render() {
        let {messages, } = this.props
        let {showScrollDown} = this.state
        let {user} = this.context
        messages = messages || []
        let events = messages.map((msg, idx) => ChatFeedEvent(msg, idx, user && msg.uid === user.uid))
        return (
            <div ref={this.feedRef} className='chatFeedContainer'>
                {showScrollDown ? (
                    <Button className='scrollDownHandler' circular icon='angle down' size='big'
                            onClick={this.scrollDown}/>
                ) : null}
                <Feed>
                    {events}
                </Feed>
            </div>
        )

    }
}

ChatFeed.contextType = StateContext

export default ChatFeed = observer(ChatFeed)
