import React from "react";
import {Button, Icon} from "semantic-ui-react";
import {Link} from "react-router-dom";
import './Sidebar.css'

export default class ButtonExampleGroupShorthand extends React.Component {
    state = {
        context: 'direct'
    }

    constructor(props) {
        super(props);
        this.state.context = props.context
    }

    onSelect = context => {
        return () => {
            this.setState({context}, () => {
                this.props.onChange && this.props.onChange(context)
            })
        }
    }

    render() {
        let {context} = this.state

        return (
            <Button.Group fluid attached='top'>
                <Link to='/direct'>
                    <Button active={context === 'direct'} onClick={this.onSelect('direct')}><Icon
                    name='user'/>Directs</Button>
                    <span/>
                </Link>
                <Link to='/channel'>
                    <Button active={context === 'channel'} onClick={this.onSelect('channel')}>
                        <Icon name='users'/>Channels
                    </Button>
                </Link>
                <Link to='/profile'>
                    <Button active={context === 'profile'} onClick={this.onSelect('profile')}>
                        <Icon name='user circle'/>Profile
                    </Button>
                </Link>
            </Button.Group>

        )
    }
}