import React, {Component} from 'react'
import {Form} from "semantic-ui-react";
import TextareaAutosize from "react-textarea-autosize";


class ChatMsgForm extends Component {
    state = {value: ""}

    constructor(props) {
        super(props);
        if (props.value) {
            this.state.value = props.value
        }
    }

    onType = (e) => {
        if (e.ctrlKey && e.code === 'Enter') {
            return this.onSend(e)
        }
        let {onType} = this.props;
        onType && onType(e.target.value);
        this.setState({value: e.target.value})
    }

    onSend = (e) => {
        let {onSend} = this.props;
        let {value} = this.state;
        let msg = value.replace(/(^\s+|\s+$)/mi, '')
        if (msg !== '' && onSend) {
            onSend(msg);
        }
        this.setState({value: ""});
    }

    render() {
        return (
            <Form >
                <Form.Group>
                    <Form.Field
                        control={TextareaAutosize}
                        placeholder="Tell us what you want..."
                        onChange={this.onType}
                        onKeyUp={this.onType}
                        value={this.state.value}
                        width={13}
                        style={{height: '2em'}}
                        rows={2}
                    />
                    <Form.Button
                        content="Отправить"
                        onClick={this.onSend}
                        width={3}
                    />
                </Form.Group>
            </Form>
        )
    }
}

export default ChatMsgForm
