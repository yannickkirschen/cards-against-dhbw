import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import './index.css';
import App from './App';
import { ThemeProvider, createTheme } from '@mui/material';

const theme = createTheme({
    palette: {
        black: {
            light: '#ffffff',
            main: '#ffffff',
            contrastText: '#ffffff',
            dark: '#ffffff',
        },
        white: {
            light: '#000000',
            main: '#000000',
            contrastText: '#000000',
            dark: '#000000',
        },
        golden: {
            light: '#0000ff',
            main: '#0000ff',
            contrastText: '#0000ff',
            dark: '#0000ff',
        }
    }
})

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <React.StrictMode>
        <ThemeProvider theme={theme}>
            <BrowserRouter>
                <App />
            </BrowserRouter>
        </ThemeProvider>
    </React.StrictMode>
);
