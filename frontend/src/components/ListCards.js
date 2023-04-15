import { Stack } from "@mui/material"
import { styled } from '@mui/material/styles';
import GameCardDisplay from "./GameCardDisplay";


const Item = styled(Paper)(({ theme }) => ({
    backgroundColor: theme.palette.mode === 'dark' ? '#1A2027' : '#fff',
    ...theme.typography.body2,
    padding: theme.spacing(1),
    textAlign: 'center',
    color: theme.palette.text.secondary,
}));

const ListCards = ({ cards, onCardSelect = null }) => {
    return (
        <Stack direction="row">
            {cards.map(el => <Item key={el}><GameCardDisplay card={el} onClick={onCardSelect} /></Item>)}
        </Stack>
    )
}

export default ListCards
