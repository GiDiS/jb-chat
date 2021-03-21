import React from 'react'
import State from "./State";
import Socket from "./Socket";

const isProd = process.env.NODE_ENV === "production"
const state = new State()
const socket = new Socket({state, server: isProd ? '' : 'ws://localhost:8888/ws'})
state.socket = socket
state.startPing(5000)

export const StateContext = React.createContext(state);

