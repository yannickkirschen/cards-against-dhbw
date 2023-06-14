import { Button } from "@mui/material";


/**
 * Component to display the current state of the game in a readable form based on the state.
 * @param {Object} props contain the state, player type, methods to determine the winner and method to leave game
 * @returns an <h3>-element containing readable information about the current state of the game
 */
const GameInfo = ({ state, isPlayerType, findWinner, leaveGame }) => {
    return (
        <>
            {(state.actionState === "game.joined") &&
                <h3>Waiting for others to join. gameCode: {localStorage.getItem("gameCode")}</h3>}

            {(state.actionState === "game.ready") &&
                <h3>Waiting for the MOD to start the game</h3>}

            {(state.actionState === "player.choosing" && isPlayerType(e => e.isBoss)) &&
                <h3>Waiting for the players to choose their cards.</h3>}

            {(state.actionState === "player.choosing" && !isPlayerType(e => e.isBoss)) &&
                <h3>Choose the card matching the black card best!</h3>}

            {(state.actionState === "boss.choosing" && isPlayerType(e => e.isBoss)) &&
                <h3>You are the boss! Choose the winning card!</h3>}

            {(state.actionState === "boss.choosing" && !isPlayerType(e => e.isBoss)) &&
                <h3>Waiting for the boss to choose the winning card</h3>}

            {(state.actionState === "mod.can.continue" && !isPlayerType(e => e.isMod)) &&
                <h3>Waiting for the MOD to continue</h3>}

            {(state.actionState === "mod.can.continue" && isPlayerType(e => e.isMod)) &&
                <h3>You are the MOD. Continue to the next round!</h3>}

            {(state.actionState === "game.finished") &&

                <h3>Game finished! {findWinner()} has won!</h3>}
            {(state.actionState === "game.abort") &&
                <h3>Game aborted. Return to home.</h3>}

            {(state.actionState === "game.error") &&
                <h3>An error ocurred: {state.msg}</h3>}

            {(state.actionState === "game.cleared") &&
                <h3>Game has been cleared: {state.msg} <br />
                    <Button variant='contained' color="error" onClick={leaveGame}>Leave Game</Button>
                </h3>}

            {state.actionState === "default" &&
                <h3>Oops...something went wrong (or I did it again)
                    Meaning either the server is still processing your request or crashed. </h3>}
        </>
    )
}

export default GameInfo;
