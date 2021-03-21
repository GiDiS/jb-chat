import React, {useContext} from 'react'
import {Button, Grid, Icon, Segment} from 'semantic-ui-react'
import {StateContext} from "../store/StateContext";
import {DirectUrl} from "../sidebar/utils";
import GoogleLogin, {GoogleLogout} from "react-google-login";
import {observer} from "mobx-react";

const responseGoogle = (response) => {
    console.log(response);
}

class LoginForm extends React.Component {
    state = {googleClientId: null, logoUrl: '/ui/logo.png'}

    render = () => {
        let {config} = this.context
        let {logoUrl} = config || {}
        console.log('login render', config)

        return (
            <Grid textAlign='center' style={{height: '100vh'}} verticalAlign='middle'>
                <Grid.Column style={{maxWidth: 450}}>
                    <Segment basic as='h2' textAlign='center'>
                        <img src={logoUrl} style={{width: '342px'}} alt=''/>
                    </Segment>
                    <Segment>
                        <Button.Group vertical>
                            <GoogleAuth/>
                            <TokenAuth/>
                        </Button.Group>
                    </Segment>
                </Grid.Column>
            </Grid>
        )
    }
}

LoginForm.contextType = StateContext

const TokenAuth = observer(() => {
    let ctx = useContext(StateContext)

    const onTyrionSignIn = () => {
        ctx.signInTyrion().then(() => {
            window.location = DirectUrl()
        })
    }

    return (
        <Button size='large' basic onClick={onTyrionSignIn}>
            <Button.Content>
                <Icon name='child'/>
                Sign as Tyrion
            </Button.Content>
        </Button>
    )
})

const GoogleAuth = observer(() => {
    let ctx = useContext(StateContext)
    let {config} = ctx
    let {googleClientId} = config || {}

    const onGoogleSignIn = (response) => {
        if (response.googleId) {
            ctx.signInGoogle({
                profileId: response.googleId,
                accessToken: response.accessToken,
                tokenId: response.tokenId
            }).then(() => {
                window.location = DirectUrl()
            })
        }
    }

    const onGoogleSignOut = (response) => {
        console.log(response)
    }

    return googleClientId ? (
        <>
            <GoogleLogin
                clientId={googleClientId}
                render={renderProps => (
                    <Button size='large' basic onClick={renderProps.onClick}>
                        <Icon name='google'/>
                        Sign in with google
                    </Button>
                )}
                buttonText="Login"
                onSuccess={onGoogleSignIn}
                onFailure={responseGoogle}
                cookiePolicy={'single_host_origin'}
            />
            <GoogleLogout
                clientId={googleClientId}
                render={renderProps => (
                    <Button size='large' basic onClick={renderProps.onClick}>
                        <Icon name='google'/>
                        Sign out with google
                    </Button>
                )}
                onLogoutSuccess={onGoogleSignOut}
            />
        </>
    ) : null

})

/*
<Form size='large'>
                    <Segment>
                        <Form.Input fluid icon='user' iconPosition='left' placeholder='E-mail address'/>
                        <Form.Input
                            fluid
                            icon='lock'
                            iconPosition='left'
                            placeholder='Password'
                            type='password'
                        />
                        <Button color='teal' fluid size='large'>
                            Login
                        </Button>
                    </Segment>
                </Form>
                <Message>
                    New to us? <a href={AccountSignUpUrl()}>Sign Up</a>
                </Message>
 */

export default observer(LoginForm)