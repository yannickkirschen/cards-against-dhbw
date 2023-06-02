# Cards against DHBW

[![Lint commit message](https://github.com/yannickkirschen/cards-against-dhbw/actions/workflows/commit-lint.yml/badge.svg)](https://github.com/yannickkirschen/cards-against-dhbw/actions/workflows/commit-lint.yml)
[![Go Workflow](https://github.com/yannickkirschen/cards-against-dhbw/actions/workflows/go.yml/badge.svg)](https://github.com/yannickkirschen/cards-against-dhbw/actions/workflows/go.yml)
[![NPM Workflow](https://github.com/yannickkirschen/cards-against-dhbw/actions/workflows/npm.yml/badge.svg)](https://github.com/yannickkirschen/cards-against-dhbw/actions/workflows/npm.yml)
[![GitHub release](https://img.shields.io/github/release/yannickkirschen/cards-against-dhbw.svg)](https://github.com/yannickkirschen/cards-against-dhbw/releases/)

This project is part of the lecture *Web Engineering 2* and is a clone of the popular card game *Cards against Humanity*.

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/yannickkirschen/cards-against-dhbw/tree/gitpod)

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
  "id": "abcdefg",
  "text": "This is a card",
  "type": 0
}
```

A card contains one and only one blank.

### Persona

- **Master of Disaster (MOD)**: The MOD is the person who creates the game and is responsible for the game flow. They can
start the game, add new cards and kick players.
- **Player**: A player can join a game and play the game. They can also add new cards to the game.

A player can be a MOD at the same time and own several white cards. They are described by the following attributes
in JSON notation:

```json
{
  "name": "John Doe",
  "active": true,
  "points": 6
}
```

### Game

A user can create a game that is identified by a unique game code which is used to join the game. Games are made up of
exact one MOD and one ore more additional players. The MOD can start the game and add new cards to the game.
