import React, {useContext} from 'react'
import {Button, Grid, Icon, Message, Segment} from 'semantic-ui-react'
import {StateContext} from "../store/StateContext";
import {DirectUrl} from "../sidebar/utils";
import GoogleLogin, {GoogleLogout} from "react-google-login";

const responseGoogle = (response) => {
    console.log(response);
}

const LoginForm = () => {
    let ctx = useContext(StateContext)

    const onTyrionSignIn = () => {
        ctx.signInTyrion().then(() => {
            window.location = DirectUrl()
        })
    }
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

    return (
        <Grid textAlign='center' style={{height: '100vh'}} verticalAlign='middle'>
            <Grid.Column style={{maxWidth: 450}}>
                <Segment basic as='h2' textAlign='center'>
                    <img src='/ui/logo.png' style={{width: '342px'}} alt=''/>
                </Segment>
                <Message>
                    <Button.Group vertical>
                        <GoogleLogin
                            clientId="20867412579-g93uhf57fabjq4q5vqk1v72pfs579dau.apps.googleusercontent.com"
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
                            clientId="20867412579-g93uhf57fabjq4q5vqk1v72pfs579dau.apps.googleusercontent.com"
                            render={renderProps => (
                                <Button size='large' basic onClick={renderProps.onClick}>
                                    <Icon name='google'/>
                                    Sign out with google
                                </Button>
                            )}
                            onLogoutSuccess={onGoogleSignOut}
                        />
                        <Button size='large' basic onClick={onTyrionSignIn}>
                            <Icon name='child'/>
                            Sign as Tyrion
                        </Button>
                    </Button.Group>
                </Message>
            </Grid.Column>
        </Grid>
    )
}

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

export default LoginForm