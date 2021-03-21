import React, {useState} from 'react'
import {Button, Feed} from 'semantic-ui-react'
import './ChatFeed.css'
import {StateContext} from '../store/StateContext'
import {ProfileLink} from "../sidebar/utils";
import {observer} from "mobx-react";
import {Link} from "react-router-dom";

const ChatFeedEvent = (msg, msgIdx, self) => {
    let {images, user} = msg;
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

    scrollToLast = () => {
        let feed = null;
        if (this.feedRef.current) {
            feed = this.feedRef.current.getElementsByClassName('feed')[0] || null
        }
        if (!feed) return false;
        let last = feed.children.length ? feed.children[feed.children.length - 1] : null;
        if (last) {
            last.scrollIntoView({block: "center", behavior: "smooth"});
            this.setLastSeen()
        }
    }

    setLastSeen = () => {
        let {messages,} = this.props
        if (messages.length) {
            let lastMsg = messages[messages.length-1]
            this.context.setChannelLastSeen(lastMsg.cid, lastMsg.mid)
        }
    }


    render() {
        let {messages,} = this.props
        let {user} = this.context
        messages = messages || []
        let events = messages.map((msg, idx) => ChatFeedEvent(msg, idx, user && msg.uid === user.uid))
        setTimeout(() => {
            if (messages.length && messages[messages.length - 1].uid === user.uid) {
                this.scrollToLast()
            }
        }, 500)
        return (
            <div ref={this.feedRef} className='chatFeedContainer'>
                <Feed>
                    {events}
                </Feed>
                <ScrollDown feedRef={this.feedRef} scrollToLast={this.scrollToLast}/>
            </div>
        )

    }
}

const ScrollDown = ({feedRef, scrollToLast}) => {
    let feed = null;
    if (feedRef.current) {
        feed = feedRef.current.getElementsByClassName('feed')[0] || null
    }
    const checkScrollTop = () => {
        if (!feed) return false;
        return feed.clientHeight < feed.scrollHeight && (feed.scrollHeight - feed.scrollTop - feed.clientHeight) > 50
    };

    const [showScrollDown, setShowScroll] = useState(checkScrollTop())
    if (!feedRef) {
        return null
    }

    const updateScrollTop = (e) => {
        setShowScroll(checkScrollTop())
    }

    const scrollDown = () => {
        scrollToLast()
        setShowScroll(false)
    }

    if (feed) {
        feed.addEventListener('scroll', updateScrollTop)
    }

    setInterval(updateScrollTop, 1000)

    let display = showScrollDown ? 'flex' : 'none';
    return (
        <Button className='scrollDownHandler' circular icon='angle down' size='big' onClick={scrollDown}
                style={{display: display}}/>
    );
}

ChatFeed.contextType = StateContext

export default ChatFeed = observer(ChatFeed)
