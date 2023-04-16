import { Route, Routes, Link } from 'react-router-dom';
import Home from './components/Home';
import GameHandler from './components/GameHandler';
import './App.css';

function App() {
    return (
        <div className="app">
            <header className="app-header">
                Cards Against DHBW - <Link to="/" >Home</Link>
            </header>
            <div className='app-body'>
                <Routes>
                    <Route path="/game/:id" element={<GameHandler />} />
                    <Route index element={<Home />} />
                </Routes>
            </div>
        </div>
    );
}

export default App;
