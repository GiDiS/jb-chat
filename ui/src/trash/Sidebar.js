
import React from 'react'
import {
    Checkbox,
    Grid,
    Header,
    Icon,
    Image,
    Menu,
    Ref,
    Segment,
    Sidebar,
} from 'semantic-ui-react'

const SidebarExampleTarget = () => {
    const segmentRef = React.useRef()
    const [visible, setVisible] = React.useState(false)

    return (
        <Grid columns={1}>
            <Grid.Column>
                <Checkbox
                    checked={visible}
                    label={{ children: <code>visible</code> }}
                    onChange={(e, data) => setVisible(data.checked)}
                />
            </Grid.Column>

            <Grid.Column>
                <Sidebar.Pushable as={Segment.Group} raised>
                    <Sidebar
                        as={Menu}
                        animation='uncover'
                        direction='left'
                        icon='labeled'
                        onHide={() => setVisible(false)}
                        vertical
                        target={segmentRef}
                        visible={visible}
                    >
                        <Menu.Item as='a' header>
                            File Permissions
                        </Menu.Item>
                        <Menu.Item as='a'>Share on Social</Menu.Item>
                        <Menu.Item as='a'>Share by E-mail</Menu.Item>
                        <Menu.Item as='a'>Edit Permissions</Menu.Item>
                        <Menu.Item as='a'>Delete Permanently</Menu.Item>
                        <Menu.Item as='a'><Icon name='home' /> Home</Menu.Item>
                        <Menu.Item as='a'><Icon name='gamepad' />Games</Menu.Item>
                        <Menu.Item as='a'><Icon name='camera' /> Channels</Menu.Item>
                    </Sidebar>

                    <Ref innerRef={segmentRef}>
                        <Segment secondary>
                            <Header as='h3'>Clickable area</Header>
                            <p>When you will click there, the sidebar will be closed.</p>
                        </Segment>
                    </Ref>

                    <Segment>
                        <Header as='h3'>Application Content</Header>
                        <Image src='https://react.semantic-ui.com/images/wireframe/paragraph.png' />
                    </Segment>
                </Sidebar.Pushable>
            </Grid.Column>
        </Grid>
    )
}

export default SidebarExampleTarget