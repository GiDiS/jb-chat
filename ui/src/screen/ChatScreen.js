import React from "react";
import {Container, Grid} from "semantic-ui-react";
import Sidebar from "../sidebar/Sidebar";
import ChatMain from "../chat/ChatMain";
import './ChatScreen.css'

export default function ChatScreen({match,location  }) {

    return (
        <Container className='App-wrapper'>
            <Grid columns='equal' className='App-greed'>
                <Grid.Row>
                    <Grid.Column width={5} className='App-sidebar'>
                        <Sidebar match={match} location={location}/>
                    </Grid.Column>
                    <Grid.Column width={11} className='App-main'>
                        <ChatMain match={match}/>
                    </Grid.Column>
                </Grid.Row>
            </Grid>
        </Container>
    )
}