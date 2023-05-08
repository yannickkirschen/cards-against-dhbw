import { Component } from "react";
import { GameCard, CardColor, Player } from "./dataStructure";

import ListPlayer from './ListPlayer'
import ListCards from './cardDisplay/ListCards'
import GameCardDisplay from "./cardDisplay/GameCardDisplay";
import { Button } from "@mui/material";

import { SocketContext } from "./socket";
import "./gameHandler.css"
import withHistory from "./hocs";
import GameInfo from "./GameInfo";

class GameHandler extends Component {

    static contextType = SocketContext;

    constructor() {
        super()
        this.state = {
            player: [],
            whiteCards: [],
            playedCards: [],
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
        this.leaveGame = this.leaveGame.bind(this);
        this.findWinner = this.findWinner.bind(this);
        this.kickPlayer = this.kickPlayer.bind(this);
        this.clearGame = this.clearGame.bind(this);
        this.isPlayerType = this.isPlayerType.bind(this);
        this.onCardDelete = this.onCardDelete.bind(this);
    }

    isPlayerType(f) {
        if (this.state.player.length === 0) {
            return false
        }

        let filtered = this.state.player.filter(f)
        if (filtered.length === 0) {
            return false
        }

        return filtered[0].name === localStorage.getItem("playerName")
    }

    loadPlayer(src) {
        let p = []
        src.forEach(element => {
            p.push(new Player(element.name, element.isMod, element.isBoss, element.points))
        });
        if (p.filter(e => e.name === localStorage.getItem("playerName")).length != 1) {
            this.leaveGame()
        }
        return p
    }

    loadWhiteCards(src) {
        let w = []
        src.forEach(el => {
            w.push(new GameCard(CardColor.WHITE, el.text, el.id))
        })
        return w
    }

    findWinner() {
        let max = 0;
        let winnerName = "";
        this.state.player.forEach(p => {
            if (p.points > max) {
                max = p.points;
                winnerName = p.name
            }
        })
        return winnerName;
    }

    onCardSelection(c) {
        this.context.emit("entity.card.chosen", JSON.stringify({ cardId: c.id }))
    }

    onCardDelete(c) {
        this.context.emit("player.card.delete", JSON.stringify({ cardId: c.id }))
    }

    leaveGame() {
        this.context.emit("game.leave", JSON.stringify({}))
        localStorage.clear()
        this.props.navigate("/")
    }

    kickPlayer(playerName) {
        this.context.emit("player.kick", JSON.stringify({ playerName: playerName }))
    }

    clearGame(reason) {
        localStorage.clear();
        this.setState({ player: [], whiteCards: [], blackCard: null, playedCards: [], actionState: reason !== null ? reason : "" })
    }

    componentDidMount() {
        this.context.emit("game.join", JSON.stringify({ gameCode: localStorage.getItem("gameCode"), playerName: localStorage.getItem("playerName") }))


        this.context.on("reconnect", attempts => {
            console.log("reconnected to server after " + attempts + " attempts")
            this.context.emit("join", JSON.stringify({ gameCode: localStorage.getItem("gameCode"), playerName: localStorage.getItem("playerName") }))
        })


        this.context.on("game.lobby", data => {
            this.setState({ player: this.loadPlayer(data.players), actionState: data.gameReady ? "game.ready" : "game.joined" })
        })

        this.context.on("player.choosing", data => {
            let b = new GameCard(CardColor.BLACK, data.blackCard.text, data.blackCard.id)
            this.setState({ player: this.loadPlayer(data.players), whiteCards: this.loadWhiteCards(data.whiteCards), blackCard: b, actionState: "player.choosing" })
        })

        this.context.on("boss.choosing", data => {
            console.log("recv: boss.choosing state")
            let b = new GameCard(CardColor.BLACK, data.blackCard.text, data.blackCard.id)
            this.setState({ blackCard: b, playedCards: this.loadWhiteCards(data.whiteCards), player: this.loadPlayer(data.players), actionState: "boss.choosing" })
        })

        this.context.on("boss.chosen", data => {
            console.log("recv: boss.chosen state")
            let cards = []
            console.log("data: " + JSON.stringify(data))
            data.playedCards.forEach((card) => {
                let nC = new GameCard(data.winnerCard === card.id ? CardColor.GOLDEN : CardColor.WHITE, card.text, card.id)
                cards.push(nC)
            })
            console.log("cards: " + JSON.stringify(cards))
            let b = new GameCard(CardColor.BLACK, data.blackCard.text, data.blackCard.id)
            this.setState({ blackCard: b, player: this.loadPlayer(data.players), playedCards: cards, actionState: "mod.can.continue" })
        })


        this.context.on("game.finished", data => {
            this.setState({ player: this.loadPlayer(data.players), actionState: "game.finished" })
            localStorage.clear()
        })

    }

    render() {
        return (
            <div className="gameHandler-container">
                <div className="gameHandler-infoArea">
                    <GameInfo state={this.state.actionState} isPlayerType={this.isPlayerType} findWinner={this.findWinner} />
                </div>
                <div className="gameHandler-blackCard">
                    {this.state.blackCard !== null && <GameCardDisplay card={this.state.blackCard} />}
                </div>
                <div className="gameHandler-publicCards">
                    <h3>Cards on the table:</h3>
                    {(this.state.actionState === "boss.choosing" || this.state.actionState === "mod.can.continue") &&
                        <div >
                            {(this.state.actionState === "boss.choosing" && this.isPlayerType(e => e.isBoss)) ?
                                <ListCards cards={this.state.playedCards} onCardSelect={this.onCardSelection} onCardDelete={this.onCardDelete} />
                                :
                                <ListCards cards={this.state.playedCards} onCardDelete={this.onCardDelete} />}
                        </div>}
                </div>
                <div className="gameHandler-whiteCards">
                    <h3>Your Cards:</h3>
                    {this.state.actionState === "player.choosing" && !this.isPlayerType(e => e.isBoss) ?
                        <ListCards cards={this.state.whiteCards} onCardSelect={this.onCardSelection} />
                        :
                        <ListCards cards={this.state.whiteCards} />}
                </div>

                {(this.isPlayerType(e => e.isMod) && this.state.actionState === "game.ready") &&
                    <Button variant="contained" color="secondary" onClick={() => this.context.emit("game.start", JSON.stringify({}))}>
                        Start game
                    </Button>}

                {(this.isPlayerType(e => e.isMod) && this.state.actionState === "mod.can.continue") &&
                    <Button variant="contained" color="secondary" onClick={() => this.context.emit("mod.round.continue", JSON.stringify({}))}>
                        Next Round
                    </Button>}

                <div className="gameHandler-playerList">
                    <h4>Player:</h4> <br />
                    <div className="gameHandler-playerList-innerContainer">
                        {!(this.isPlayerType(e => e.isMod)) ? <ListPlayer player={this.state.player} /> : <ListPlayer player={this.state.player} kickHandler={this.kickPlayer} />}
                    </div>
                    <Button variant='contained' color="error" onClick={this.leaveGame}>Leave Game</Button>
                </div>
                <div className="gameHandler-devInfo">
                    playerName: {localStorage.getItem("playerName")},
                    gameCode: {localStorage.getItem("gameCode")},
                    Action State: {this.state.actionState},
                    IsMod: {this.isPlayerType(e => e.isMod) ? "yes" : "no"},
                    IsBoss: {this.isPlayerType(e => e.isBoss) ? "yes" : "no"}
                </div>
            </div>

        )
    }

}

export default withHistory(GameHandler);
