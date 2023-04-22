import { Paper, Stack } from "@mui/material"
import { styled } from '@mui/material/styles';
import GameCardDisplay from "./GameCardDisplay";
import GameCardButton from "./GameCardButton";
import { CardColor } from "../dataStructure";


const Item = styled(Paper)(({ theme }) => ({
    backgroundColor: theme.palette.mode === 'dark' ? '#1A2027' : '#fff',
    ...theme.typography.body2,
    padding: theme.spacing(1),
    textAlign: 'center',
    color: theme.palette.text.secondary,
}));

const ListCards = ({ cards, onCardSelect = null }) => {
    console.log("displaying cards: " + JSON.stringify(cards))
    return (
        <Stack direction="row" flexWrap="wrap" display="flex" justifyContent={"space-evenly"} margin={"auto"}>
            {cards.map(el =>
                <Item key={JSON.stringify(el)} sx={{ width: el.color === CardColor.BLACK ? "80%" : "16%", display: "flex", flexDirection: "column", justifyContent: "space-between", marginBottom: "10px" }}>
                    <GameCardDisplay card={el} />

                    <GameCardButton card={el} onCardSelect={onCardSelect} />
                </Item>)}
        </Stack >
    )
}

export default ListCards
