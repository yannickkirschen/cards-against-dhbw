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

        this.socket.on("lobbyState", data => {
            console.log("recv: lobby state")
            this.setState({ player: this.loadPlayer(data.players), actionState: "gameWaiting" })
        })

        this.socket.on("playerChoosingState", data => {
            console.log("recv: playerChoosingState state")
            let b = new GameCard(CardColor.BLACK, data.blackCard.text, data.blackCard.id)
            this.setState({ player: this.loadPlayer(data.players), whiteCards: this.loadWhiteCards(data.whiteCards), blackCard: b, actionState: "playerChoosing" })
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

        this.socket.on("bossChoosingState", data => {
            console.log("recv: bossChoosingState state")
            let b = new GameCard(CardColor.BLACK, data.blackCard.text, data.blackCard.id)
            this.setState({ blackCard: b, whiteCards: this.loadWhiteCards(data.playedCards), player: this.loadPlayer(data.players), actionState: "bossChoosing" })
        })

        this.socket.on("bossHasChosenState", data => {
            console.log("recv: boss has chosen state")
            let cards = []
            data.playedCards.forEach((owner, card, map) => {
                let nC = new GameCard(data.winnerCard.text.id == card.id ? CardColor.GOLDEN : CardColor.WHITE, card.text, card.id)
                nC.setOwner(owner)
                cards.push(nC)
            })
            let b = new GameCard(CardColor.BLACK, data.blackCard.text, data.blackCard.id)
            this.setState({ blackCard: b, player: this.loadPlayer(data.players), whiteCards: cards, actionState: "none" })
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
                    <div className="gameHandler-whiteCards">
                        {this.state.actionState == "playerChoosing" || this.state.actionState == "bossChoosing" ?
                            <ListCards cards={this.state.whiteCards} onCardSelect={this.onCardSelection} />
                            :
                            <ListCards cards={this.state.whiteCards} />}
                        {console.log("state: " + JSON.stringify(this.state))}
                    </div>
                    <div className="gameHandler-blackCard">
                        <ListCards cards={[this.state.blackCard]} />
                    </div>
                </>
            )
        }
        else {
            return (
                <>
                    <Button onClick={() => { this.socket.emit("notice", "msgTOll") }} >Click Me</Button>
                    <Button onClick={() => { this.socket.emit("joinRequestAction", JSON.stringify({ gameID: "gello", playerID: "jello" })) }} >Click Me</Button>
                    <h2>Connection to game server failed.</h2>
                    <Link to="/">Home</Link>
                </>
            )
        }
    }
}
