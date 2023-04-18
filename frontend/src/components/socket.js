import * as io from 'socket.io-client'

const socket = io.connect("http://192.168.0.26:3333/", { transports: ['websocket', 'polling'] })


export default socket;
