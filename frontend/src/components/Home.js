import { useNavigate, useParams } from 'react-router-dom'
import { Button, Divider, TextField } from '@mui/material'
import { useEffect, useState } from 'react'

import { newGame, joinGame } from './functions'

import './Home.css'

/**
 * Component that is shown at the beginning.
 * offers inputs for the player to either join a game or create a new one
 * @returns component that contains the home-view
 */
const Home = () => {
    const navigate = useNavigate();
    const urlID = useParams().id;
    const [gameCode, setGameCode] = useState("")
    const [joinName, setJoinName] = useState("")

    const [newName, setNewName] = useState("")

    useEffect(() => {
        console.log("urlID: " + urlID);
        setGameCode(urlID !== undefined ? urlID : "")
    }, [])

    return (
        <div className='home'>
            <div className='home-join-game'>
                <h2> Join Game </h2>
                <TextField
                    fullWidth
                    label="gameCode"
                    variant="filled"
                    value={gameCode}
                    onChange={(e) => setGameCode(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && gameCode && joinName ? joinGame(gameCode, joinName, navigate) : null} />
                <br />
                <TextField
                    fullWidth
                    label="Name"
                    variant="filled"
                    value={joinName}
                    onChange={(e) => setJoinName(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && gameCode && joinName ? joinGame(gameCode, joinName, navigate) : null} />
                <br />
                <Button variant='contained' color='primary' disabled={!gameCode || !joinName} onClick={() => joinGame(gameCode, joinName, navigate)}>Join Game</Button>
            </div>
            <Divider flexItem />
            <div className='home-create-game'>
                <h2> Create Game </h2>
                <TextField
                    fullWidth
                    label="Name"
                    variant="filled"
                    value={newName}
                    onChange={(e) => setNewName(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && newName ? newGame(newName, navigate) : null} />
                <br />
                <Button variant='contained' color='primary' disabled={!newName} onClick={() => newGame(newName, navigate)}>Create New Game</Button>
            </div>
        </div>
    )
}

export default Home;
