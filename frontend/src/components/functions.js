
const newGame = async (gamerName, navigate) => {
    fetch("http://localhost:3333/v1/new?name=" + gamerName)
        .then(res => res.json())
        .then(data => {
            console.log("data: " + JSON.stringify(data))
            localStorage.setItem("gameID", data.gameId)
            localStorage.setItem("playerID", data.playerId)
            return data.gameId
        })
        .then(gameId => {
            navigate("/game/" + gameId)
        })
        .catch(err => {
            console.log("caught error while fetching new game: " + err)
            localStorage.setItem("gameID", "emptyID")
            localStorage.setItem("playerID", "emptyPlayerID")
            navigate("/game/emptyID")
        })
}

const joinGame = async (gameID, gamerName, navigate) => {
    fetch("http://localhost:3333/v1/join/" + gameID + "?name=" + gamerName)
        .then(res => res.json())
        .then(data => {
            localStorage.setItem("gameID", data.gameId)
            localStorage.setItem("playerID", data.playerId)
            return data.gameId
        })
        .then(gameId => {
            navigate("/game/" + gameID.Id)
        })
        .catch(err => {
            console.log("caught error while fetching new game: " + err)
            localStorage.setItem("gameID", "emptyID")
            localStorage.setItem("playerID", "emptyPlayerID")
            navigate("/game/emptyID")
        })
}


export { newGame, joinGame };
