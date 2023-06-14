import { Route, Routes } from 'react-router-dom';
import Home from './components/Home';
import GameHandler from './components/GameHandler';
import './App.css';
import { socket, SocketContext } from './components/socket';
import Header from './components/header/Header';
import HelpPage from './components/help/Help';

function App() {
    return (
        <SocketContext.Provider value={socket}>
            <div className="app">
                <Header />
                <div className='app-body'>
                    <Routes>
                        <Route path="/game" element={<GameHandler />} />
                        <Route path="/help" element={<HelpPage />} />
                        <Route path="/:id?" element={<Home />} />

                    </Routes>
                </div>
            </div>
        </SocketContext.Provider>
    );
}

export default App;
