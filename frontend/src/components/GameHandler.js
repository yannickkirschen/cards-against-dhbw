import { Component } from "react";
import * as io from 'socket.io-client'
import { GameCard, CardColor, Player } from "./dataStructure";

import ListPlayer from './ListPlayer'
import ListCards from './ListCards'
import { Link } from "react-router-dom";
import { Button } from "@mui/material";


export default class GameHandler extends Component {

    constructor() {
        super()
        this.state = {
            player: [],
            whiteCards: [],
            blackCard: new GameCard(CardColor.BLACK, "", "NONE"),
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

    isPlayerBoss() {
        return this.state.player.filter(e => e.isBoss)[0].id == localStorage.getItem("playerID")
    }

    loadPlayer(src) {
        let p = []
        src.forEach(element => {
            p.push(new Player(element.name, element.isMod, element.isBoss, element.points))
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
        this.socket = io.connect("http://localhost:8000", { transports: ['websocket', 'polling'] })
        this.socket.on("connect", () => {
            console.log("recv: connect")
            this.socket.emit("joinRequestAction", JSON.stringify({ gameID: localStorage.getItem("gameID"), playerID: localStorage.getItem("playerID") }))

        })
        this.didUnMount = () => { this.socket.disconnect(); console.log("disconnected ") }

        this.socket.on("invalidCodeState", data => {
            console.log("recv: invalid code state")
            this.socket.disconnect()
            this.setState({ actionState: "invalidCode" })
            localStorage.removeItem("gameID")
            localStorage.removeItem("playerID")
        })

        this.socket.on("game.joined", data => {
            console.log("recv: lobby state")
            this.setState({ player: this.loadPlayer(data.players), actionState: "game.joined" })
        })

        this.socket.on("game.finished", data => {
            console.log("recv: lobby state")
            this.setState({ player: this.loadPlayer(data.players), actionState: "game.finished" })
        })


        this.socket.on("player.choosing", data => {
            console.log("recv: player.choosing")
            let b = new GameCard(CardColor.BLACK, data.blackCard.text, data.blackCard.id)
            this.setState({ player: this.loadPlayer(data.players), whiteCards: this.loadWhiteCards(data.whiteCards), blackCard: b, actionState: "player.choosing" })
        })

        this.socket.on("bossWaitingState", data => {
            console.log("recv: boss waiting state")
            let b = new GameCard(CardColor.BLACK, data.blackCard.text, data.blackCard.id)
            let w = []
            for (let card = 0; card < data.numberPlayedCards; card++) {
                w.push(new GameCard(CardColor.WHITE, "Hidden Card", "NONE"))
            }
            this.setState({ player: this.loadPlayer(data.players), whiteCards: w, blackCard: b, actionState: "none" })
        })

        this.socket.on("boss.choosing", data => {
            console.log("recv: boss.choosing state")
            let b = new GameCard(CardColor.BLACK, data.blackCard.text, data.blackCard.id)
            this.setState({ blackCard: b, playedCards: this.loadWhiteCards(data.whiteCards), player: this.loadPlayer(data.players), actionState: "boss.choosing" })
        })

        this.socket.on("boss.chosen", data => {
            console.log("recv: boss.chosen state")
            let cards = []
            data.playedCards.forEach((card) => {
                let nC = new GameCard(data.winnerCard == card.id ? CardColor.GOLDEN : CardColor.WHITE, card.text, card.id)
                cards.push(nC)
            })
            let b = new GameCard(CardColor.BLACK, data.blackCard.text, data.blackCard.id)
            this.setState({ blackCard: b, player: this.loadPlayer(data.players), playedCards: cards, actionState: "boss.can.continue" })
        })
    }

    onCardSelection(c) {
        this.socket.emit("cardChosenAction", JSON.stringify({ card: c.id }))
    }

    componentWillUnmount() {
        console.log("component will unmount called")
        this.didUnMount()
    }

    render() {
        if (this.state.actionState !== "invalidCode") {
            return (
                <>
                    <div className="gameHandler-playerList">
                        <ListPlayer player={this.state.player} />
                    </div>
                    {(this.state.actionState == "boss.choosing" || this.state.actionState == "boss.can.continue") && <div className="gameHandler-playedCards">
                        {(this.state.actionState == "boss.choosing" && this.isPlayerBoss()) ?
                            <ListCards cards={this.state.playedCards} onCardSelect={this.onCardSelection} />
                            :
                            <ListCards cards={this.state.playedCards} />}
                    </div>}
                    <div className="gameHandler-whiteCards">
                        {this.state.actionState == "playerChoosing" && !this.isPlayerBoss() ?
                            <ListCards cards={this.state.whiteCards} onCardSelect={this.onCardSelection} />
                            :
                            <ListCards cards={this.state.whiteCards} />}
                    </div>
                    <div className="gameHandler-blackCard">
                        <ListCards cards={[this.state.blackCard]} />
                    </div>
                    {(this.isPlayerBoss() && this.state.actionState == "boss.can.continue") && <Button variant="contained" onClick={() => this.so
                        .emit("game.start", JSON.stringify({}))}>Next Round</Button>}
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
