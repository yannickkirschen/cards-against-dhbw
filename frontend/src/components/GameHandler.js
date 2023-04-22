import { Component } from "react";
import * as io from 'socket.io-client'
import { GameCard, CardColor, Player } from "./dataStructure";

import ListPlayer from './ListPlayer'
import ListCards from './ListCards'
import { Link } from "react-router-dom";
import { Button } from "@mui/material";

import socket from "./socket";

export default class GameHandler extends Component {

    constructor() {
        super()
        this.state = {
            player: [],
            whiteCards: [],
            blackCard: null,
            actionState: "invalidCoe",
            // action states are:
            //  invalidCode - link to home enabled
            //  playerChoosing - whiteCard buttons enabled
            //  bossChoosing - whiteCard buttons enabled
            //  gameWaiting - startGame button enabled
            //  none
        }
        this.onCardSelection = this.onCardSelection.bind(this);

    }

    isPlayerType(f) {
        console.log(this.state.player)
        if (this.state.player.length === 0) {
            return false
        }

        let filtered = this.state.player.filter(f)
        console.log(filtered)
        if (filtered.length === 0) {
            return false
        }

        return filtered[0].id === localStorage.getItem("playerID")
    }

    loadPlayer(src) {
        let p = []
        console.log("loadPlayer called")
        src.forEach(element => {
            p.push(new Player(element.id, element.name, element.isMod, element.isBoss, element.points))
        });
        return p
    }

    loadWhiteCards(src) {
        let w = []
        src.forEach(el => {
            w.push(new GameCard(CardColor.WHITE, el.text, el.id))
        })
        return w
    }

    componentDidMount() {
        console.log("component did Mount called")
        // socket = io.connect("http://192.168.0.26:3333/", { transports: ['websocket', 'polling'] })
        if (socket.connected) {
            socket.emit("join", JSON.stringify({ gameID: localStorage.getItem("gameID"), playerID: localStorage.getItem("playerID") }))
        }
        this.didUnMount = () => { socket.disconnect(); console.log("disconnected ") }

        socket.on("invalidCodeState", data => {
            console.log("recv: invalid code state")
            this.setState({ actionState: "invalidCode" })
            localStorage.removeItem("gameID")
            localStorage.removeItem("playerID")
        })

        socket.on("game.joined", data => {
            console.log("recv: lobby state  > " + JSON.stringify(data))
            this.setState({ player: this.loadPlayer(data.players), actionState: data.gameReady ? "game.ready" : "game.joined" })
        })

        socket.on("game.finished", data => {
            console.log("recv: lobby state")
            this.setState({ player: this.loadPlayer(data.players), actionState: "game.finished" })
        })


        socket.on("player.choosing", data => {
            console.log("recv: player.choosing")
            console.log(JSON.stringify(data))
            let b = new GameCard(CardColor.BLACK, data.blackCard.text, data.blackCard.id)
            this.setState({ player: this.loadPlayer(data.players), whiteCards: this.loadWhiteCards(data.whiteCards), blackCard: b, actionState: "player.choosing" })
        })

        socket.on("boss.choosing", data => {
            console.log("recv: boss.choosing state")
            let b = new GameCard(CardColor.BLACK, data.blackCard.text, data.blackCard.id)
            this.setState({ blackCard: b, playedCards: this.loadWhiteCards(data.whiteCards), player: this.loadPlayer(data.players), actionState: "boss.choosing" })
        })

        socket.on("boss.chosen", data => {
            console.log("recv: boss.chosen state")
            let cards = []
            data.playedCards.forEach((card) => {
                let nC = new GameCard(data.winnerCard === card.id ? CardColor.GOLDEN : CardColor.WHITE, card.text, card.id)
                cards.push(nC)
            })
            let b = new GameCard(CardColor.BLACK, data.blackCard.text, data.blackCard.id)
            this.setState({ blackCard: b, player: this.loadPlayer(data.players), playedCards: cards, actionState: "mod.can.continue" })
        })
    }

    onCardSelection(c) {
        socket.emit("cardChosenAction", JSON.stringify({ card: c.id }))
    }

    componentWillUnmount() {
        console.log("component will unmount called")
        this.didUnMount()
    }

    render() {
        if (this.state.actionState !== "invalidCode") {
            console.log(this.isPlayerType(e => e.isMod))
            return (
                <>
                    <h3>PlayerID: {localStorage.getItem("playerID")}, GameID: {localStorage.getItem("gameID")}, Action State: {this.state.actionState}, IsMod: {this.isPlayerType(e => e.isMod) ? "yes" : "no"}, IsBoss: {this.isPlayerType(e => e.isBoss) ? "yes" : "no"}</h3>
                    <div className="gameHandler-playerList">
                        <ListPlayer player={this.state.player} />
                    </div>
                    {(this.state.actionState === "boss.choosing" || this.state.actionState === "boss.can.continue") && <div className="gameHandler-playedCards">
                        {(this.state.actionState === "boss.choosing" && this.isPlayerType(e => e.isBoss)) ?
                            <ListCards cards={this.state.playedCards} onCardSelect={this.onCardSelection} />
                            :
                            <ListCards cards={this.state.playedCards} />}
                    </div>}
                    <div className="gameHandler-whiteCards">
                        {this.state.actionState === "player.choosing" && !this.isPlayerType(e => e.isBoss) ?
                            <ListCards cards={this.state.whiteCards} onCardSelect={this.onCardSelection} />
                            :
                            <ListCards cards={this.state.whiteCards} />}
                    </div>
                    <div className="gameHandler-blackCard">
                        {this.state.blackCard !== null && <ListCards cards={[this.state.blackCard]} />}
                    </div>
                    {(this.isPlayerType(e => e.isMod) && this.state.actionState === "game.ready") && <Button variant="contained" onClick={() => socket.emit("game.start", JSON.stringify({}))}>Start game</Button>}

                    {(this.isPlayerType(e => e.isMod) && this.state.actionState === "mod.can.continue") && <Button variant="contained" onClick={() => socket.emit("mod.round.continue", JSON.stringify({}))}>Next Round</Button>}
                </>
            )
        }
        else {
            return (
                <>
                    <h2>Connection to game server failed.</h2>
                    <Link to="/">Home</Link>
                </>
            )
        }
    }
}
