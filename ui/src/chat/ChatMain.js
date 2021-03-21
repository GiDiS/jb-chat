import React, {Component} from "react";
import Feed from "./ChatFeed";
import ChatMsgForm from "./ChatMsgForm";
import './ChatMain.less'
import {Button, Dimmer, DimmerDimmable, Header, Icon, Loader, Segment} from "semantic-ui-react";
import {observer} from "mobx-react";
import {StateContext} from "../store/StateContext";

class ChatMain extends Component {
    state = {loading: false, creatingRoom: false, roomTitle: null};

    componentDidMount() {
        let activeChannelId = this.context.activeChannelId || null;
        if (activeChannelId) {
            this.context.getActiveChannel(activeChannelId).then((d) => {
                console.log(d)
            })
        }
    }

    onMessage = msg => {
        this.context.addMessage(msg)
    }

    onLeave = () => {
        this.context.leaveChannel()
    }

    onJoin = () => {
        this.context.joinChannel()
    }

    onDrop = () => {
        this.context.dropChannel()
    }

    onLikeClick = (msg) => {
        console.log(msg)
        this.context.likeMessage(msg.cid, msg.mid)
    }

    channelActions = (chan, user) => {
        if (!chan) {
            return null
        }
        let actions = [];

        const isPublic = chan.isPublic && chan.isPublic()
        const isMember = chan.isMember && chan.isMember(user);
        if (isPublic && chan.cid && isMember) {
            actions.push((<Button size='tiny' floated='right' onClick={this.onLeave} key='chan-leave'>Leave</Button>))
        }
        if (isPublic && chan.cid && !isMember) {
            actions.push((<Button size='tiny' floated='right' onClick={this.onJoin} key='chan-join'>Join</Button>))
        }

        if (isMember) {
            actions.push((<Button size='tiny' floated='right' onClick={this.onDrop} key='chan-join'>Drop</Button>))
        }

        return actions
    }

    render() {
        let activeChannelId = this.context.activeChannelId || null;
        let activeChannel = this.context.activeChannel || {};
        let chatLoading = this.context.chatLoading || false;
        let user = this.context.user;

        // console.log(activeChannelId)
        // console.log(activeChannel)

        const messages = activeChannel.messages || [];
        const members = activeChannel.users || [];
        const membersStatus = this.context.getUsersStatus(members)
        const isDirect = activeChannel.cid && activeChannel.isDirect()
        const isWriteable = activeChannel.cid && activeChannel.isWritable(user, membersStatus)
        const title = activeChannel.title || 'Channel not selected';
        let icon = isDirect ? 'user' : 'users'
        return (
            <div className='ChatMain'>
                <Segment.Group>
                    <Segment>
                        <Header>
                            <Header.Content as='strong'>
                                <Icon name={icon}/>
                                {title}
                                <small> (cid: {activeChannel.cid || activeChannelId || ""}, members: {members.length},
                                    messages: {messages.length})</small>
                                {this.channelActions(activeChannel, user)}
                            </Header.Content>
                        </Header>

                    </Segment>
                    <Segment>
                        <DimmerDimmable className='flex-vert'>
                            <Dimmer active={chatLoading} inverted>
                                <Loader size='large'/>
                            </Dimmer>
                            <Feed messages={messages} onLikeClick={this.onLikeClick}/>
                        </DimmerDimmable>
                    </Segment>
                    {isWriteable ? (
                        <Segment>
                            <ChatMsgForm onSend={this.onMessage}/>
                        </Segment>
                    ) : null}
                </Segment.Group>
            </div>
        )
    }
}


ChatMain.contextType = StateContext

export default ChatMain = observer(ChatMain)