import { useNavigate } from "react-router-dom";
import HomeIcon from '@mui/icons-material/Home';
import HelpIcon from '@mui/icons-material/Help';
import "./header.css"
import { IconButton } from "@mui/material";

function Header() {
    const navigate = useNavigate()
    return (
        <div>
            <header className="App-header">
                <IconButton sx={{ borderRadius: '5px' }} className="homeIcon-button" onClick={() => navigate("/")}>
                    <HomeIcon color="secondary" sx={{ fontSize: 50 }} className='homeIcon-icon' />
                </IconButton>
                <div className='header-text'>
                    Cards Against DHBW
                </div>
                <IconButton className="helpIcon-button" sx={{ borderRadius: "5px" }} onClick={() => navigate("/help")}>
                    <HelpIcon color="secondary" sx={{ fontSize: 50 }} className="helpIcon-icon" />
                </IconButton>

            </header>
        </div>
    )
}

export default Header;
