import { List, ListItemIcon, ListItemText } from "@mui/material"
import AccountBoxIcon from '@mui/icons-material/AccountBox';
import StarIcon from '@mui/icons-material/Star';

const ListPlayer = ({ player }) => {
    return (
        <List>
            {player.map(p =>
                <ListItem key={p.name}>
                    {p.isBoss ? <ListItemIcon> <StarIcon /></ListItemIcon> : <ListItemIcon><AccountBoxIcon /> </ListItemIcon>}
                    <ListItemText primary={p.name + (p.isMod ? " (MOD)" : "")} secondary={"Points:" + p.points} />
                </ListItem>
            )}
        </List>
    )
}


export default ListPlayer
