import React from 'react';
import ReactDOM from 'react-dom';
import Dashboard from './js/Dashboard';
import Menu from './js/Menu';
import {loadAuth} from './js/auth'
import './scss/app.scss';

const menuEl = document.getElementById('menu');
const el = document.getElementById('dashboard');

loadAuth();

ReactDOM.render(<Menu />, menuEl);

if (el) {
  ReactDOM.render(<Dashboard wshost={el.getAttribute('data-wshost')} />, el);
}
