import { ButtonBase, Card, CardContent } from "@mui/material";
import { CardColor } from "./dataStructure";

const GameCardDisplay = ({ card, onCardSelect = null }) => {
    let backColor = card.color;
    let textColor = card.color == CardColor.BLACK ? "white" : "black"
    if (onCardSelect == null) {
        return (
            <Card color="primary">
                <CardContent style={{ backgroundColor: backColor, color: textColor }}>
                    {card.content}
                </CardContent>
            </Card>
        )
    }
    else {
        return (
            <Card style={{ backgroundColor: backColor, color: textColor }}>
                <ButtonBase onClick={() => onCardSelect(card)}>
                    <CardContent>
                        {card.content}
                    </CardContent>
                </ButtonBase>
            </Card>
        )
    }
}

export default GameCardDisplay;
