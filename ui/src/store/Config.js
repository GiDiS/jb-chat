import {makeAutoObservable} from 'mobx'

export default class Config {
    logoUrl = '/ui/logo.png'
    googleClientId = null

    update = (props = {}) => {
        if (props.logo_url) this.logoUrl = props.logo_url
        if (props.logoUrl) this.logoUrl = props.logoUrl
        if (props.google_client_id) this.googleClientId = props.google_client_id
        if (props.googleClientId) this.googleClientId = props.googleClientId
    }

    constructor(props) {
        this.update(props)
        makeAutoObservable(this)
    }
}