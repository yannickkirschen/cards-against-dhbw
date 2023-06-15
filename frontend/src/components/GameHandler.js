import { Component } from "react";
import { GameCard, CardColor, Player } from "./dataStructure";
import ErrorSnackbar from "./ErrorSnackbar";

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

    /**
     * constructor method.
     * Initiates the state-variable and binds the context to the functions.
     */
    constructor() {
        super()
        this.state = {
            player: [],
            whiteCards: [],
            playedCards: [],
            blackCard: null,
            actionState: "default",
            showError: false,
            msg: ""
        }
        this.onCardSelection = this.onCardSelection.bind(this);
        this.leaveGame = this.leaveGame.bind(this);
        this.findWinner = this.findWinner.bind(this);
        this.kickPlayer = this.kickPlayer.bind(this);
        this.isPlayerType = this.isPlayerType.bind(this);
        this.onCardDelete = this.onCardDelete.bind(this);
    }

    /**
     * checks if the current player is of a certain type (either mod or boss)
    * @param    {Function}  f   the function that is used to filter the player array
    * @returns  {Boolean}       if the player is of the type that is tested by the parameter f
    */
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

    /**
     * helper function that parses a list of players into an array of Player-objects
     * @param {Array} src   the array containing the players to be loaded
     * @returns             array of player-objects
     */
    loadPlayer(src) {
        let p = []
        src.forEach(element => {
            p.push(new Player(element.name, element.isMod, element.isBoss, element.points, element.active))
        });
        if (p.filter(e => e.name === localStorage.getItem("playerName")).length != 1) {
            this.leaveGame()
        }
        return p
    }

    /**
     * helper function that parses a list of whiteCards into an array of GameCard-objects
     * @param {Array} src   the array containing the cards to be loaded
     * @returns             an array containing GameCard-objects
     */
    loadWhiteCards(src) {
        let w = []
        src.forEach(el => {
            w.push(new GameCard(CardColor.WHITE, el.text, el.id))
        })
        return w
    }

    /**
     * function to find the player with the most points in the game.
     * @returns the name of the current player with the most points, "" if the game has not started
     */
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

    /**
     * function to send a message to the server that the player wants to select a card
     * @param {Card} c the card that has been selected
     */
    onCardSelection(c) {
        this.context.emit("entity.card.chosen", JSON.stringify({ cardId: c.id }))
    }

    /**
     * function to send a message to the server that the player wants to remove a card
     * @param {Card} c the card that has been selected to be removed
     */
    onCardDelete(c) {
        this.context.emit("player.card.remove", JSON.stringify({ cardId: c.id }))
    }

    /**
     * function to send a message to the server that the player wants to leave the game
     */
    leaveGame() {
        this.context.emit("game.leave", JSON.stringify({}))
        localStorage.clear()
        this.props.navigate("/")
    }

    /**
     * function to send a message to the server that the player wants to kick a player.
     * This fails if the player emitting the event is NOT the mod.
     * @param {String} playerName the name of the player that is supposed to be kicked
     */
    kickPlayer(playerName) {
        this.context.emit("player.kick", JSON.stringify({ playerName: playerName }))
    }

    /**
     * function that is called when the GameHandler-component is mounted.
     * emits a join-event to the server
     * sets up the listeners for events sent by the server via the socket connection
     */
    componentDidMount() {
        this.context.emit("game.join", JSON.stringify({ gameCode: localStorage.getItem("gameCode"), playerName: localStorage.getItem("playerName") }))

        this.context.on("reconnect", attempts => {
            console.log("reconnected to server after " + attempts + " attempts")
            this.context.emit("game.rejoin", JSON.stringify({ gameCode: localStorage.getItem("gameCode"), playerName: localStorage.getItem("playerName") }))
        })

        /**
         * on the game.lobby event the players in the game are loaded into the state and the gameState is updated.
         */
        this.context.on("game.lobby", data => {
            this.setState({ player: this.loadPlayer(data.players), actionState: data.gameReady ? "game.ready" : "game.joined" })
        })

        /**
         * on the player.choosing event the players in the game are loaded into the state as well as the white cards and the black card. The action state is set to player.choosing.
         * If the player is not the boss, he is now allowed to choose a card.
         */
        this.context.on("player.choosing", data => {
            let b = new GameCard(CardColor.BLACK, data.blackCard.text, data.blackCard.id)
            this.setState({ player: this.loadPlayer(data.players), whiteCards: this.loadWhiteCards(data.whiteCards), blackCard: b, actionState: "player.choosing" })
        })

        /**
         * on the boss.choosing state event the players in the game, black card and the white cards are loaded into the state.
         * If the player is the boss, he is now allowed to choose a card.
         */
        this.context.on("boss.choosing", data => {
            console.log("recv: boss.choosing state")
            let b = new GameCard(CardColor.BLACK, data.blackCard.text, data.blackCard.id)
            this.setState({ blackCard: b, playedCards: this.loadWhiteCards(data.whiteCards), player: this.loadPlayer(data.players), actionState: "boss.choosing" })
        })

        /**
         * on the boss.chosen state event the players in the game, the black card, the white cards and the played cards are loaded into the state.
         * The winning card is saved with a different card than the other cards.
         * If the player is the mod, he can now continue the game to the next round.
         */
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

        /**
         * on the game.finished state event the players in the game are loaded into the state.
         * The local storage is cleared. The players will see the winner.
         */
        this.context.on("game.finished", data => {
            this.setState({ player: this.loadPlayer(data.players), actionState: "game.finished" })
            localStorage.clear()
        })

        /**
         * on the game.error state event an error message is displayed to the player.
         */
        this.context.on("game.error", data => {
            this.setState({ actionState: "game.error", msg: data.label + ": " + data.payload, showError: true })
        })
    }

    render() {
        return (
            <div className="gameHandler-container">
                <div className="gameHandler-infoArea">
                    <GameInfo state={this.state} isPlayerType={this.isPlayerType} findWinner={this.findWinner} leaveGame={this.leaveGame} />
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
                        <ListCards cards={this.state.whiteCards} onCardSelect={this.onCardSelection} onCardDelete={this.onCardDelete} />
                        :
                        <ListCards cards={this.state.whiteCards} onCardDelete={this.onCardDelete} />}
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
                    <h4>GameID: {localStorage.getItem("gameCode")} <br /> <br /> Player:</h4> <br /> <br />
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
                <ErrorSnackbar msg={this.state.msg} open={this.state.showError} setOpen={(val) => this.setState({ showError: val })} />
            </div>
        )
    }
}

export default withHistory(GameHandler);
