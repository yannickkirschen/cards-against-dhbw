import { Card, CardContent } from "@mui/material";

const GameCardDisplay = ({ card, onCardSelect = null }) => {
    if (onCardSelect == null) {
        return (
            <Card color={card.color}>
                <CardContent>
                    {card.content}
                </CardContent>
            </Card>
        )
    }
    else {
        return (
            <Card color={card.color} onClick={() => onCardSelect(card)}>
                <CardContent>
                    {card.content}
                </CardContent>
            </Card>
        )
    }
}

export default GameCardDisplay;
