import React, {Component} from "react";
import {Form, Icon, Label} from "semantic-ui-react";

export default class Filter extends Component {
    constructor(props) {
        super(props);
        this.inputRef = React.createRef()
    }
    onChange = event => {
        const {updateFilter} = this.props
        updateFilter(event.target.value)
    }

    onReset = event => {
        const {updateFilter} = this.props
        updateFilter("")
        event.target.value = ""; console.log(this.inputRef)
    }

    render() {

        return (
            <div style={{padding: '0 1rem'}}>
                <Form.Input
                    labelPosition='right'
                    type='text'
                    placeholder='Search...'
                    size='mini'
                >
                    <Label basic><Icon name='search' /></Label>
                    <input ref={this.inputRef} onChange={this.onChange} size='mini' style={{width: "100%"}}/>
                    <Label><Icon name='delete' onClick={this.onReset}/></Label>
                </Form.Input>
            </div>
        )
    }
}