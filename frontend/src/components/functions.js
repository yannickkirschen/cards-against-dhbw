
const newGame = async (gamerName, navigate) => {
    fetch("/v1/new?name=" + gamerName)
        .then(res => res.json())
        .then(data => {
            localStorage.setItem("gameID", data.id)
            localStorage.setItem("playerID", data.playerId)
        })
        .then(data => {
            navigate("/game/" + data.id)
        })
        .catch(err => {
            console.log("caught error while fetching new game: " + err)
            localStorage.setItem("gameID", "emptyID")
            localStorage.setItem("playerID", "emptyPlayerID")
            navigate("/game/emptyID")
        })
}

const joinGame = async (gameID, gamerName, navigate) => {
    fetch("/v1/join/" + gameID + "?name=" + gamerName)
        .then(res => res.json())
        .then(data => {
            localStorage.setItem("gameID", data.id)
            localStorage.setItem("playerID", data.playerId)
        })
        .then(data => {
            navigate("/game/" + data.id)
        })
        .catch(err => {
            console.log("caught error while fetching new game: " + err)
            localStorage.setItem("gameID", "emptyID")
            localStorage.setItem("playerID", "emptyPlayerID")
            navigate("/game/emptyID")
        })
}


export { newGame, joinGame };
