import * as io from 'socket.io-client'
import React from 'react'

export const socket = io.connect("/", { transports: ['websocket', 'polling'] })
export const SocketContext = React.createContext();
