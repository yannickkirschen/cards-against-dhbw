import { Component } from "react";
import { GameCard, CardColor, Player } from "./dataStructure";

import ListPlayer from './ListPlayer'
import ListCards from './cardDisplay/ListCards'
import GameCardDisplay from "./cardDisplay/GameCardDisplay";
import { Link } from "react-router-dom";
import { Button } from "@mui/material";

import { SocketContext } from "./socket";
import "./gameHandler.css"
import withHistory from "./hocs";

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

    }

    isPlayerType(f) {
        if (this.state.player.length === 0) {
            return false
        }

        let filtered = this.state.player.filter(f)
        if (filtered.length === 0) {
            return false
        }

        return filtered[0].id === localStorage.getItem("playerID")
    }

    loadPlayer(src) {
        let p = []
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
        this.context.emit("cardChosenAction", JSON.stringify({ card: c.id }))
    }

    leaveGame() {
        this.context.emit("game.leave", JSON.stringify({}));
        this.props.navigate("/")
    }

    componentWillUnmount() {
        this.leaveGame();
    }

    componentDidMount() {
        if (this.context.connected) {
            this.context.emit("join", JSON.stringify({ gameID: localStorage.getItem("gameID"), playerID: localStorage.getItem("playerID") }))
        }
        if (!this.context.connected) {
            // somehow this happens on page reload, but connection is
            this.context.emit("join", JSON.stringify({ gameID: localStorage.getItem("gameID"), playerID: localStorage.getItem("playerID") }))
        }

        // socket = io.connect("http://192.168.0.26:3333/", { transports: ['websocket', 'polling'] })
        // const socket = io.connect("http://localhost:3333/", { transports: ['websocket', 'polling'] })
        //this.didUnMount = () => { this.context.disconnect(); console.log("disconnected ") }

        this.context.on("invalidCodeState", data => {
            this.setState({ actionState: "invalidCode" })
            localStorage.removeItem("gameID")
            localStorage.removeItem("playerID")
        })

        this.context.on("game.joined", data => {
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
        })

    }

    componentWillUnmount() {
        console.log("component will unmount called")
        //this.didUnMount()
    }

    render() {
        if (this.state.actionState !== "invalidCode") {
            console.log(this.isPlayerType(e => e.isMod))
            return (
                <div className="gameHandler-container">
                    <div className="gameHandler-infoArea">
                        {(this.state.actionState === "game.joined") &&
                            <h3>Waiting for others to join. GameID: {localStorage.getItem("gameID")}</h3>}
                        {(this.state.actionState === "game.ready") &&
                            <h3>Waiting for the MOD to start the game</h3>}
                        {(this.state.actionState === "player.choosing" && this.isPlayerType(e => e.isBoss)) &&
                            <h3>Waiting for the players to choose their cards.</h3>}
                        {(this.state.actionState === "player.choosing" && !this.isPlayerType(e => e.isBoss)) &&
                            <h3>Choose the card matching the black card best!</h3>}
                        {(this.state.actionState === "boss.choosing" && this.isPlayerType(e => e.isBoss)) &&
                            <h3>You are the boss! Choose the winning card!</h3>}
                        {(this.state.actionState === "boss.choosing" && !this.isPlayerType(e => e.isBoss)) &&
                            <h3>Waiting for the boss to choose the winning card</h3>}
                        {(this.state.actionState === "mod.can.continue" && !this.isPlayerType(e => e.isMod)) &&
                            <h3>Waiting for the MOD to continue</h3>}
                        {(this.state.actionState === "mod.can.continue" && this.isPlayerType(e => e.isMod)) &&
                            <h3>You are the MOD. Continue to the next round!</h3>}
                        {(this.state.actionState === "game.finished") &&
                            <h3>Game finished! {this.findWinner()} has won!</h3>}
                        {this.state.actionState === "invalidCoe" &&
                            <h3>Oops...something went wrong (or I did it again)</h3>}


                    </div>
                    <div className="gameHandler-blackCard">
                        {this.state.blackCard !== null && <GameCardDisplay card={this.state.blackCard} />}
                    </div>
                    <div className="gameHandler-publicCards">
                        <h3>Cards on the table:</h3>
                        {(this.state.actionState === "boss.choosing" || this.state.actionState === "mod.can.continue") &&
                            <div >
                                {(this.state.actionState === "boss.choosing" && this.isPlayerType(e => e.isBoss)) ?
                                    <ListCards cards={this.state.playedCards} onCardSelect={this.onCardSelection} />
                                    :
                                    <ListCards cards={this.state.playedCards} />}
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
                            <ListPlayer player={this.state.player} />
                        </div>
                        <Button variant='contained' color="error" onClick={this.leaveGame}>Leave Game</Button>
                    </div>
                    <div className="gameHandler-devInfo">
                        PlayerID: {localStorage.getItem("playerID")},
                        GameID: {localStorage.getItem("gameID")},
                        Action State: {this.state.actionState},
                        IsMod: {this.isPlayerType(e => e.isMod) ? "yes" : "no"},
                        IsBoss: {this.isPlayerType(e => e.isBoss) ? "yes" : "no"}
                    </div>
                </div>
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

export default withHistory(GameHandler);
