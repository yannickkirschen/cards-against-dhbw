import * as io from 'socket.io-client'

const socket = io.connect("http://localhost:3333/", { transports: ['websocket', 'polling'] })


export default socket;
