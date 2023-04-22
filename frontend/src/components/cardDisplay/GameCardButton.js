import { ButtonBase, Card, CardContent } from "@mui/material";

const GameCardButton = ({ card, onCardSelect }) => {
    if (onCardSelect !== null) {
        return (
            <Card style={{ clear: "both", marginTop: '10px' }} >
                <ButtonBase onClick={() => onCardSelect(card)}>
                    <CardContent>
                        Select this card
                    </CardContent>
                </ButtonBase>
            </Card>
        )
    }
    return <></>
}

export default GameCardButton;
