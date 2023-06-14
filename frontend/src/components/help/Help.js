import "./help.css"


const HelpPage = () => {
    return (
        <div className="help-body">
            <h2>How To: Cards Against DHBW</h2>
            Anybody can go to the homepage and create a new game by entering an ingame-name, click the "Create Game"-button and share the code to their friends. <br /> <br />
            After everybody has joined the game (at least two players are required), <br />the MOD (master of disaster, i.e. the person that created the game) can start the game by pressing the "Start Game"-button. <br />
            Each player then receives a random & unique set of 10 white cards from the card-pool and one black card is displayed to all players. <br />
            Every player, except the Boss, now has to select one of their white cards of which they believe it matches the blank in the black card the best. <br /> <br />
            After everyone, except the Boss, has selected a white card, the action jumps to the Boss, who has to select the funniest white card (with respect to the black card presented). <br />
            The player, whose white card was selected, receives one point, all players' white cards are refilled and the next round begins, with another player being the Boss. <br />
            If one player receives 10 points, they win the game. <br />
            <h3>Have Fun!!</h3>
        </div>
    )
}

export default HelpPage;
