import { Card, CardContent } from "@mui/material";
import { CardColor } from "../dataStructure";

const GameCardDisplay = ({ card }) => {
    let backColor = card.color;
    let textColor = card.color === CardColor.BLACK ? "white" : "black"

    return (
        <Card color="primary">
            <CardContent style={{ backgroundColor: backColor, color: textColor }}>
                {card.content}
            </CardContent>
        </Card>
    )
}

export default GameCardDisplay;
