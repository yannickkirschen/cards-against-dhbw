import * as io from 'socket.io-client'
import React from 'react'

export const socket = io.connect("http://localhost:3333/", { transports: ['websocket', 'polling'] })
export const SocketContext = React.createContext();
