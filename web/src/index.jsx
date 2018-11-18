import React from 'react';
import ReactDOM from 'react-dom';
import Dashboard from './js/Dashboard';
import Menu from './js/Menu';
import ThemeLoader from './js/ThemeLoader';
import ThemeProvider from './js/ThemeProvider';
import './scss/app.scss';

const menuEl = document.getElementById('menu');
const dashboard = document.getElementById('dashboard');

const App = (
  <ThemeProvider>
    <ThemeLoader />
    <Menu />
    {dashboard && <Dashboard  wshost={dashboard.getAttribute('data-wshost')} />}
  </ThemeProvider>
);

ReactDOM.render(App, document.getElementById('app'));
