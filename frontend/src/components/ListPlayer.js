import { List, ListItem, ListItemIcon, ListItemText, Button } from "@mui/material"
import AccountBoxIcon from '@mui/icons-material/AccountBox';
import StarIcon from '@mui/icons-material/Star';

const ListPlayer = ({ player, kickHandler }) => {
    console.log("kick handler: " + kickHandler)
    return (
        <List>
            {player.map(p =>
                <ListItem key={p.name}>
                    {p.isBoss ? <ListItemIcon> <StarIcon /></ListItemIcon> : <ListItemIcon><AccountBoxIcon /> </ListItemIcon>}
                    <ListItemText primary={p.name + (p.isMod ? " (MOD)" : "")} secondary={"Points:" + p.points} />
                    {kickHandler !== undefined && <Button color="error" onClick={() => kickHandler(p.name)}>Kick</Button>}
                </ListItem>
            )}
        </List>
    )
}




export default ListPlayer
