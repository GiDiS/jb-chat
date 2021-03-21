import React, {Component} from "react"
import {observer} from "mobx-react";
import {Button, Form, Icon, Item, Modal} from "semantic-ui-react";
import {StateContext} from '../store/StateContext';
import Filter from "./Filter";
import './ChannelsList.css'
import {GotoChannel} from "./utils";
import {ChannelStatus} from "./ChannelStatus";

class ChannelsList extends Component {
    state = {filter: '', showCreateDialog: false}

    componentDidMount() {
        let {getChannels} = this.context
        getChannels()
    }

    updateFilter = (filter) => {
        this.setState({filter})
    }

    onSelectChannel = (chan) => {
        return (e) => {
            this.context.setActiveChannelId(chan.cid)
            GotoChannel(chan)(e)
        }
    }

    render() {
        let {channels, activeChannelId} = this.context
        let {filter} = this.state
        channels = channels.filter(chan => chan.isPublic())

        if (filter) {
            filter = filter.toLowerCase()
            channels = channels.filter(ch => ch.title.toLowerCase().includes(filter))
        }
        let items = channels.map((chan, idx) => {
            let isSelected = activeChannelId && activeChannelId === chan.cid
            let className = 'channelItem'
            if (isSelected) {
                className += ' active'
            }
            return (
                <Item key={'chan_' + idx} className={className} onClick={this.onSelectChannel(chan)}>
                    {/*<Link to={ChannelLink(chan)} onClick={this.onSelectChannel}>*/}
                    <Item.Content verticalAlign='middle' className='channelItem-name'>
                        {chan.title}
                        <ChannelStatus channel={chan}/>
                    </Item.Content>
                    {/*</Link>*/}
                </Item>
            )
        })

        return (
            <Item.Group className='channelsList'>
                <Filter updateFilter={this.updateFilter}/>
                <ModalCreateChannel/>

                {items}
            </Item.Group>

        )
    }
}

ChannelsList.contextType = StateContext


class ModalCreateChannel extends React.Component {
    state = {value: "", error: "", open: false}

    onClose = () => {
        this.setState({open: false})
    }

    onOpen = () => {
        this.setState({open: true})
    }

    onDone = () => {
        if (this.state.value) {
            this.context.createChannel({title: this.state.value})
            this.setState({open: false, value: ""})
        }
    }

    render() {
        let {open} = this.state

        return (
            <Modal
                open={open}
                onClose={() => this.onClose}
                onOpen={() => this.onOpen}
                trigger={
                    <Item key={'chan_new'} className='channelItem' onClick={this.onOpen}>
                        <Item.Content verticalAlign='middle' className='channelItem-name'>
                            <span><Icon name='plus square outline'/>Create channel</span>
                        </Item.Content>
                    </Item>
                }
            >
                <Modal.Header>Create channel!</Modal.Header>
                <Modal.Content>
                    <Form.Input
                        text
                        className='channelName'
                        onChange={event => this.setState({value: event.target.value})}
                        placeholder='Please enter channel name...'
                    />
                </Modal.Content>
                <Modal.Actions>
                    <Button onClick={this.onClose}>Close</Button>
                    <Button onClick={this.onDone} positive>Create</Button>
                </Modal.Actions>
            </Modal>
        )
    }
}

ModalCreateChannel.contextType = StateContext

export default ChannelsList = observer(ChannelsList)


