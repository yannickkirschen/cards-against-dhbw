


const GameInfo = ({ state, isPlayerType, findWinner }) => {
    return (
        <>
            {(state === "game.joined") &&
                <h3>Waiting for others to join. gameCode: {localStorage.getItem("gameCode")}</h3>}

            {(state === "game.ready") &&
                <h3>Waiting for the MOD to start the game</h3>}

            {(state === "player.choosing" && isPlayerType(e => e.isBoss)) &&
                <h3>Waiting for the players to choose their cards.</h3>}

            {(state === "player.choosing" && !isPlayerType(e => e.isBoss)) &&
                <h3>Choose the card matching the black card best!</h3>}

            {(state === "boss.choosing" && isPlayerType(e => e.isBoss)) &&
                <h3>You are the boss! Choose the winning card!</h3>}

            {(state === "boss.choosing" && !isPlayerType(e => e.isBoss)) &&
                <h3>Waiting for the boss to choose the winning card</h3>}

            {(state === "mod.can.continue" && !isPlayerType(e => e.isMod)) &&
                <h3>Waiting for the MOD to continue</h3>}

            {(state === "mod.can.continue" && isPlayerType(e => e.isMod)) &&
                <h3>You are the MOD. Continue to the next round!</h3>}

            {(state === "game.finished") &&

                <h3>Game finished! {findWinner()} has won!</h3>}
            {(state === "game.abort") &&
                <h3>Game aborted. Return to home.</h3>}

            {state === "invalidCoe" &&
                <h3>Oops...something went wrong (or I did it again)</h3>}
        </>
    )
}

export default GameInfo;
