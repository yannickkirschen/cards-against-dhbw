
const newGame = async (navigate) => {
    fetch("/v1/new")
        .then(res => res.json())
        .then(data => {
            localStorage.setItem("id", data.id)
            localStorage.setItem("playerId", data.playerId)
        })
        .then(data => {
            navigate("/" + data.id)
        })
        .catch(err => {
            console.log("caught error while fetching new game: " + err)
            localStorage.setItem("gameID", "emptyID")
            localStorage.setItem("playerID", "emptyPlayerID")
            navigate("/emptyID")
        })
}

const joinGame = async (gameID, navigate) => {
    fetch("/v1/join/" + gameID)
        .then(res => res.json())
        .then(data => {
            localStorage.setItem("gameID", data.id)
            localStorage.setItem("playerID", data.playerId)
        })
        .then(data => {
            navigate("/" + data.id)
        })
        .catch(err => {
            console.log("caught error while fetching new game: " + err)
            localStorage.setItem("gameID", "emptyID")
            localStorage.setItem("playerID", "emptyPlayerID")
            navigate("/emptyID")
        })
}


export { newGame, joinGame };
