# Cards against DHBW

This project is part of the lecture *Web Engineering 2* and is a clone of the popular card game *Cards against Humanity*.

## Developers

- [Rex2002](https://github.com/Rex2002)
- [yannickkirschen](https://github.com/yannickkirschen)
- [MalteRichert](https://github.com/MalteRichert)

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

- **Master of Desaster (MOD)**: The MOD is the person who creates the game and is responsible for the game flow. They can
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
