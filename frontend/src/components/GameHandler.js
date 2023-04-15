import { Component } from "react";
import * as io from 'socket.io-client'
import { GameCard, CardColor, Player } from "./dataStructure";

import ListPlayer from './ListPlayer'
import ListCards from './ListCards'
import { Link } from "react-router-dom";


export default class GameHandler extends Component {

    state = {
        player: [],
        whiteCards: [],
        blackCard: new GameCard(CardColor.BLACK, "", "NONE"),
        actionState: ""
        // action states are:
        //  invalidCode - link to home enabled
        //  playerChoosing - whiteCard buttons enabled
        //  bossChoosing - whiteCard buttons enabled
        //  gameWaiting - startGame button enabled
        //  none
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
        data.whiteCards.forEach(el => {
            w.push(new GameCard(Color.WHITE, el.text, el.id))
        })
        return w
    }

    componentDidMount() {
        this.socket = io.connect("randomIP")
        if (!this.socket.connected) {
            this.setState({ actionState: "invalidCode" })
        }
        this.didUnMount = () => { socket.disconnect() }
        socket.on("invalidCodeState", data => {
            socket.disconnect()
            this.setState({ actionState: "invalidCode" })
            localStorage.removeItem("gameID")
            localStorage.removeItem("playerID")
        })
        socket.on("lobbyState", data => {
            this.setState({ player: this.loadPlayer(data.players), actionState: "gameWaiting" })
        })
        socket.on("playerChoosingState", data => {
            let b = new GameCard(Color.BLACK, data.blackCard.text, data.blackCard.id)
            this.setState({ player: this.loadPlayer(data.players), whiteCards: this.loadWhiteCards(data.whiteCards), blackCard: b, actionState: "playerChoosing" })
        })
        socket.on("bossWaitingState", data => {
            let b = new GameCard(Color.BLACK, data.blackCard.text, data.blackCard.id)
            let w = []
            for (let card = 0; card < data.numberPlayedCards; card++) {
                w.push(new GameCard(Color.WHITE, "Hidden Card", "NONE"))
            }
            this.setState({ player: this.loadPlayer(data.players), whiteCards: w, blackCard: b, actionState: "none" })
        })
        socket.on("bossChoosingState", data => {
            let b = new GameCard(Color.BLACK, data.blackCard.text, data.blackCard.id)
            this.setState({ blackCard: b, whiteCards: this.loadWhiteCards(data.playedCards), player: this.loadPlayer(data.players), actionState: "bossChoosing" })
        })
        socket.on("bossHasChosenState", data => {
            let cards = []
            data.playedCards.forEach((owner, card, map) => {
                let nC = new GameCard(data.winnerCard.text.id == card.id ? Color.GOLDEN : Color.WHITE, card.text, card.id)
                nC.setOwner(owner)
                cards.push(nC)
            })
            let b = new GameCard(Color.BLACK, data.blackCard.text, data.blackCard.id)

            this.setState({ blackCard: b, player: this.loadPlayer(data.players), whiteCards: cards, actionState: "none" })
        })
        this.socket.emit("joinReq", { gameID: localStorage.getItem("gameID"), playerID: localStorage.getItem("playerID") })
    }

    onCardSelection(c) {
        console.log("card selected")
        this.socket.emit("cardPlayed", { card: c.id })
    }

    componentWillUnmount() {
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
                    </div>
                    <div className="gameHandler-blackCard">
                        <ListCards card={[this.state.BlackCard]} />
                    </div>
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
