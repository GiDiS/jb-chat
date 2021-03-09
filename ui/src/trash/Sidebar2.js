import React from 'react'
import {
    Button,
    Checkbox,
    Grid,
    Header,
    Icon,
    Image, Item,
    Menu, Reveal,
    Segment,
    Sidebar,
} from 'semantic-ui-react'
import '../sidebar/UsersList.css'


const RevealExampleMoveDown = () => (
    <Reveal animated='move down'>
        <Reveal.Content visible>
            <Image src='https://react.semantic-ui.com/images/wireframe/square-image.png' size='tiny'
                   style={{borderRadius: '50%'}}/>
        </Reveal.Content>
        <Reveal.Content hidden>
            <Image src='https://i.pravatar.cc/150?u=a042581f4e29026704s' size='tiny' style={{borderRadius: '50%'}}/>
        </Reveal.Content>
    </Reveal>
)



const SidebarExampleMultiple = () => {
    let [visible, setVisible] = React.useState(false)

    visible = true

    return (
        <Grid columns={1}>
            <Grid.Column>
                <Checkbox
                    checked={visible}
                    label={{children: <code>visible</code>}}
                    onChange={(e, data) => setVisible(data.checked)}
                />
            </Grid.Column>

            <Grid.Column>
                <Sidebar.Pushable as={Segment}>
                    <Sidebar
                        as={Menu}
                        animation='overlay'
                        direction='left'
                        onHide={() => setVisible(false)}
                        vertical
                        visible={visible}
                        width='wide'
                    >


                        <ButtonExampleGroupShorthand/>
                        <Menu.Item as='a'>
                            <RevealExampleMoveDown/>
                        </Menu.Item>
                        <Menu.Item as='a'>
                            <Icon name='home'/>
                            Home
                        </Menu.Item>
                        <Menu.Item as='a'>
                            <Icon name='gamepad'/>
                            Games
                        </Menu.Item>
                        <Menu.Item as='a'>
                            <Icon name='camera'/>
                            Channels
                        </Menu.Item>
                    </Sidebar>

                    <Sidebar
                        as={Menu}
                        animation='overlay'
                        direction='right'
                        vertical
                        visible={visible}
                    >
                        <Menu.Item as='a' header>
                            File Permissions
                        </Menu.Item>
                        <Menu.Item as='a'>Share on Social</Menu.Item>
                        <Menu.Item as='a'>Share by E-mail</Menu.Item>
                        <Menu.Item as='a'>Edit Permissions</Menu.Item>
                        <Menu.Item as='a'>Delete Permanently</Menu.Item>
                    </Sidebar>

                    <Sidebar.Pusher>
                        <Segment basic>
                            <Header as='h3'>Application Content</Header>
                            <Image src='https://react.semantic-ui.com/images/wireframe/paragraph.png'/>
                        </Segment>
                    </Sidebar.Pusher>
                </Sidebar.Pushable>
            </Grid.Column>
        </Grid>
    )
}

export default SidebarExampleMultiple
