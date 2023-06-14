import { List, ListItem, ListItemIcon, ListItemText, Button } from "@mui/material"
import AccountBoxIcon from '@mui/icons-material/AccountBox';
import StarIcon from '@mui/icons-material/Star';

/**
 * Component to display a list of the current players as well as their state (active/inactive).
 * If the player is the mod, the kickHandler is passed and a button for kicking a player is added.
 * @param {Object} props contain the players and a function that is called when a player is kicked
 * @returns a list containing a visualization of the current players.
 */
const ListPlayer = ({ player, kickHandler }) => {
    console.log("player: " + JSON.stringify(player))
    return (
        <List>
            {player.map(p =>
                <ListItem key={p.name}>
                    {p.isBoss ? <ListItemIcon> <StarIcon /></ListItemIcon> : <ListItemIcon> <AccountBoxIcon /> </ListItemIcon>}
                    <ListItemText primary={p.name + (p.isMod ? " (MOD)" : "") + (p.active ? "" : "(INACTIVE)")} secondary={"Points:" + p.points} />
                    {kickHandler !== undefined && <Button color="error" onClick={() => kickHandler(p.name)}>Kick</Button>}
                </ListItem>
            )}
        </List>
    )
}




export default ListPlayer
