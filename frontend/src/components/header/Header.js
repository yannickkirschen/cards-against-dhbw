import { useNavigate } from "react-router-dom";
import HomeIcon from '@mui/icons-material/Home';
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
                <div className="App-header-placeholder">

                </div>

            </header>
        </div>
    )
}

export default Header;
