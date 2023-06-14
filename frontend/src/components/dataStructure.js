/**
 * provides classes to hold data about card and players.
 */


const CardColor = {
    BLACK: "black",
    WHITE: "white",
    GOLDEN: "golden"
}

class Player {
    constructor(name, isMod, isBoss, points, active) {
        this.name = name
        this.isMod = isMod
        this.isBoss = isBoss
        this.points = points
        this.active = active
    }
}

class GameCard {
    constructor(color, content, id) {
        this.color = color
        this.id = id
        this.content = content
        this.owner = ''
    }

    setOwner(owner) {
        this.owner = owner;
    }

}

export { CardColor, Player, GameCard }
