import { Route, Routes } from 'react-router-dom';
import Home from './components/Home';
import GameHandler from './components/GameHandler';
import './App.css';

function App() {
    return (
        <div className="app">
            <header className="app-header">
                Cards Against DHBW
            </header>
            <div className='app-body'>
                <Routes>
                    <Route path="/">
                        <Route path=":id" element={<GameHandler />} />
                        <Route index element={<Home />} />
                    </Route>
                </Routes>
            </div>
        </div>
    );
}

export default App;
