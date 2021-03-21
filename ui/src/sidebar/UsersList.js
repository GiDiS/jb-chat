import React, {Component} from 'react'
import {observer} from "mobx-react";
import {Item} from "semantic-ui-react";
import "./UsersList.css"
import Filter from "./Filter";
import {GotoDirect} from "./utils";
import {StateContext} from "../store/StateContext";

class UsersList extends Component {
    state = {filter: ''}

    componentDidMount() {
        let {getUsers} = this.context
        getUsers()
    }

    updateFilter = (filter) => {
        this.setState({filter})
    }

    onChannelSelect = (user) => {
        return (e) => {
            this.context.setActiveChannelId('@' + user.nickname)
            GotoDirect(user)(e)
        }
    }

    prepareUsers = (users, filter) => {
        // remove users without avatars
        users = users.filter(u => u.avatarUrl !== "")
        // apply filter
        if (filter) {
            filter = filter.toLowerCase()
            users = users.filter(u => u.title.toLowerCase().includes(filter))
        }
        users.sort((a, b) => {
            let _a = (a.status === 'online' ? "0" : "1") + a.title.toLowerCase()
            let _b = (b.status === 'online' ? "0" : "1")  + b.title.toLowerCase()
            return _a < _b ? -1 : (_a > _b ? 1 : 0);
        })
        return users
    }

    render() {
        let {users: usersMap, activeChannelId} = this.context || {users: {}}
        let {filter} = this.state
        let users = this.prepareUsers(Object.values(usersMap), filter)
        let isSelected = user => '@' + user.nickname === activeChannelId

        let items = users.map(user => {
            let className = 'userItem'
            if (isSelected(user)) {
                className += ' active'
            }
            // let direct = directs[user.nickname] || null

            return (
                <Item key={'user_' + user.uid} className={className} onClick={this.onChannelSelect(user)}>
                    <Item.Image size='tiny' src={user.avatarUrl} className='userItem-avatar' alt={user.nickname}/>
                    <Item.Content verticalAlign='middle' className='userItem-name'>
                        {user.title}<br/>
                        <small>[{user.status}]</small>
                    </Item.Content>
                </Item>
            )
        })

        return (
            <Item.Group className='usersList'>
                <Filter updateFilter={this.updateFilter}/>
                {items}
            </Item.Group>
        )
    }
}

UsersList.contextType = StateContext

export default UsersList = observer(UsersList)