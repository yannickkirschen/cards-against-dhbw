
/**
 * a function that makes an api-request to create a new game.
 * after the game is created, forwarding the playing-endpoint in the frontend is initiated
 * @param {String} gamerName the in-game-name of the player that wants to create a game
 * @param {*} navigate a function to forward to the playing page after all necessary api-calls have been made
 */
const newGame = async (gamerName, navigate) => {
    fetch("/v1/new?name=" + gamerName)
        .then(res => res.json())
        .then(data => {
            console.log("data: " + JSON.stringify(data))
            localStorage.setItem("gameCode", data.gameCode)
            localStorage.setItem("playerName", data.playerName)
            return data.gameCode
        })
        .then(() => {
            navigate("/game")
        })
        .catch(err => {
            console.log("caught error while fetching new game: " + err)
            localStorage.setItem("gameCode", "emptyID")
            localStorage.setItem("playerName", "emptyplayerName")
            navigate("/game")
        })
}

/**
 * function to join an existing game using a game code.
 * @param {String} gameCode gameCode of the game that the player wants to join
 * @param {String} gamerName in-game-name of the player that wants to join a game
 * @param {*} navigate a function to forward to the playing page after all necessary api-calls have been made
 */
const joinGame = async (gameCode, gamerName, navigate) => {
    fetch("/v1/join/" + gameCode + "?name=" + gamerName)
        .then(res => res.json())
        .then(data => {
            localStorage.setItem("gameCode", data.gameCode)
            localStorage.setItem("playerName", data.playerName)
            return data.gameCode
        })
        .then(() => {
            navigate("/game")
        })
        .catch(err => {
            console.log("caught error while fetching new game: " + err)
            localStorage.setItem("gameCode", "emptyID")
            localStorage.setItem("playerName", "emptyplayerName")
            navigate("/game")
        })
}


export { newGame, joinGame };
