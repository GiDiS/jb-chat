

function ChannelStatus({channel}) {
    if (!channel) {
        return null
    }

    // @todo implement chan status
    if (channel.lastMessage) {
        return null
    }
    return null
}

export {
    ChannelStatus
}