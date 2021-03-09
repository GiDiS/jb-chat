import React from "react";
import {Button, Card, Icon, Image} from "semantic-ui-react";
import {observer} from "mobx-react";
import {StateContext} from "../store/StateContext";
import {GotoDirect} from "./utils";


class Profile extends React.Component {

    componentDidMount() {
        this.context.getUsers().then()
    }

    onSignOut = () => {
        this.context.signOut()
    }

    onChannelSelect = (user) => {
        return (e) => {
            this.context.setActiveChannelId('@'+user.nickname)
            GotoDirect(user)(e)
        }
    }

    render() {
        let {location} = this.props
        let {user} = this.context
        user = user || {}
        let currentUser = (user ? user.nickname : false) || false
        let match = location.pathname.match(/^\/profile\/(\w+)/)

        if (match) {
            user = this.context.getUser(match[1])
            currentUser = false
        }
        return (
            <Card fluid>
                {user ? (<UserCard user={user}/>) : null}
                {currentUser ? (<Button onClick={this.onSignOut}><Icon name='sign out'/> Sign out</Button>) : null}
                {user && !currentUser ? (<Button onClick={this.onChannelSelect(user)}><Icon name='chat'/> Direct</Button>) : null}
            </Card>
        )
    }

}

Profile.contextType = StateContext

function UserCard({user}) {
    return (
        <>
            <Image src={user.avatarUrl} wrapped ui={false}/>
            <Card.Content>
                <Card.Header>{user.title}</Card.Header>
                <Card.Meta>
                    <span className='user'>@{user.nickname}</span><br/>
                    <span className=''><a href={'mailto:'+user.email}>{user.email}</a></span><br/>
                    {/*<span className='date'>@todo Joined in 2015</span>*/}
                </Card.Meta>
                <Card.Description>
                    No description
                </Card.Description>
            </Card.Content>
            {/*<Card.Content extra>*/}
            {/*    <a>*/}
            {/*        <Icon name='user'/>*/}
            {/*        22 Friends*/}
            {/*    </a>*/}
            {/*</Card.Content>*/}
        </>
    )
}

export default Profile = observer(Profile)