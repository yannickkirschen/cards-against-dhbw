import { useNavigate } from 'react-router-dom'
import { Button, Divider, TextField } from '@mui/material'
import { useState } from 'react'

import { newGame, joinGame } from './functions'

import './Home.css'


const Home = () => {
    const navigate = useNavigate();

    const [gameID, setGameID] = useState("")
    const [joinName, setJoinName] = useState("")
    const [newName, setNewName] = useState("")

    return (
        <div className='home'>
            <div className='home-join-game'>
                <h2> Join Game </h2>
                <TextField fullWidth label="GameID" variant="filled" value={gameID} onChange={(e) => setGameID(e.target.value)} />
                <br />
                <TextField fullWidth label="Name" variant="filled" value={joinName} onChange={(e) => setJoinName(e.target.value)} />
                <br />
                <Button variant='contained' color='primary' disabled={!gameID || !joinName} onClick={() => joinGame(gameID, joinName, navigate)}>Join Game</Button>
            </div>
            <Divider flexItem />
            <div className='home-create-game'>
                <h2> Create Game </h2>
                <TextField fullWidth label="Name" variant="filled" value={newName} onChange={(e) => setNewName(e.target.value)} />
                <br />
                <Button variant='contained' color='primary' disabled={!newName} onClick={() => newGame(newName, navigate)}>Create New Game</Button>
            </div>
        </div>
    )
}

export default Home;
