import { ButtonBase, Card } from "@mui/material";
import CloseIcon from '@mui/icons-material/Close';

const GameCardDelete = ({ card, onCardDelete }) => {
    if (onCardDelete !== null) {
        return (
            <Card style={{ marginBottom: '5px', marginTop: '10px', width: '10%', }} >
                <ButtonBase onClick={() => onCardDelete(card)}>
                    <CloseIcon />
                </ButtonBase>
            </Card>
        )
    }
    return <></>
}

export default GameCardDelete;
