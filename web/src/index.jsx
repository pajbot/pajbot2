import React from 'react';
import ReactDOM from 'react-dom';
import Dashboard from './js/Dashboard';
import './scss/app.scss';

const el = document.getElementById('root');

ReactDOM.render(<Dashboard wshost={el.getAttribute('data-wshost')} />, el);