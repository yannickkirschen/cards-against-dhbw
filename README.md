# Cards against DHBW

This project is part of the lecture *Web Engineering 2* and is a clone of the popular card game *Cards against Humanity*.

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/yannickkirschen/cards-against-dhbw)

## Developers

- [Rex2002](https://github.com/Rex2002)
- [yannickkirschen](https://github.com/yannickkirschen)
- [MalteRichert](https://github.com/MalteRichert)

## Features

- Multiplayer web application for playing the game *Cards against DHBW*.
- The Master of Disaster (MOD) can create a game and start it. They receive a unique game code that can be shared with other players.
- Players can join a game by entering the game code.
- Players can choose a unique name.
- The MOD can start the game after players have joined the lobby.
- When a game is started, nobody can join the game anymore.
- Every player holds ten white cards.
- If a player doesn't like their cards, they can remove them individually before the round ends and receive new ones the next round.
- Each round presents a black card with a question or a sentence with blanks.
- Players except the Boss can play white cards to fill the blanks.
- Once every player chose their cards, the Boss chooses the funniest card to win the round.
- The player who won the round gets a point.
- The game ends when a player reaches ten points or the MOD closes the game or all players leave the game.

## Data Model

### Card

A card is a piece of text that can be used in the game. It is described by the following attributes in JSON notation:

```json
{
  "id": 1,
  "text": "This is a card",
  "type": "black/white",
  "blanks": 1
}
```

If the card is a white card, the blanks are always `-1`. If the card is a black card, the blanks are always `0` or greater.

### Persona

- **Master of Disaster (MOD)**: The MOD is the person who creates the game and is responsible for the game flow. They can
start the game, add new cards and kick players.
- **Player**: A player can join a game and play the game. They can also add new cards to the game.

A player can be a MOD at the same time and own several white cards. They are described by the following attributes
in JSON notation:

```json
{
  "id": 1,
  "name": "John Doe",
  "isMod": true,
  "cards": [1, 2, 3]
}
```

### Game

A user can create a game that is identified by a unique game code which is used to join the game. Games are made up of
exact one MOD and one ore more additional players. The MOD can start the game and add new cards to the game. The game
is described by the following attributes in JSON notation:

```json
{
  "id": 1,
  "code": "ABC123",
  "mod": 1,
  "players": [1, 2, 3],
  "blackCards": [1, 2, 3],
  "whiteCards": [1, 2, 3],
  "state": "waiting/playing/finished"
}
```

## Server

### API

- /join/:id   {name: String}
- /leave/:id
- /kill/:game {auth: ?}
  auth is required to ensure that the user is allowed to kill the game (i.e. is mod)
- /start/:id {auth: ?}

### Socket

Game logic is handled by the server.
The client only needs to know his role (cardTsar yes/no), his hand, the card in the center, the current status (e.g. able to play a card, waiting for others to play their cards) as well as his current options.
Each request should come with a gameID to allow handling of multiple games simultaneously.
From server to all clients of one game:

- killConfirm()
- updateGameState()
  type:
  - handUpdate: changes cards of the player
  - tableCardUpate: changes black card displayed
  - invisibleHandsUpdate: updates the display, if the other players have played their cards (without showing the cards)
  - gameState: able to play a card, waiting for others, cardTsar, etc. (maybe find a better name)
  - scoreBoard: update points, etc.
- updateActions()
  type:
  - playBlackCard
  - playWhiteCard
  - chooseWinningCard
  - toggleShuffleAvailability
- updateReqAckknowledged()
- ping()
  to check if all players are still connected
- playerUpdate()
  if a visual aid is implemented to indicate the number of players in the current round, it should be updated following this broadcast

From clients to server:

- killRequest()
- updateGameStateReq()
  type: play card, choose winnerCard, shuffleCards, skipNextRoundTimer (while viewing the game after the Card Tsar chose the winning card)
  value: played card
- pingAck()
