
const newGame = async (gamerName, navigate) => {
    fetch("http://localhost:3333/v1/new?name=" + gamerName)
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

const joinGame = async (gameCode, gamerName, navigate) => {
    fetch("http://localhost:3333/v1/join/" + gameCode + "?name=" + gamerName)
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
