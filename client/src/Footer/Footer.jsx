import React from 'react';
import './Footer.css';
import svgIcons from '../svgIcons';
import Navigation from './Navigation/Navigation';

export default function Footer() {
    return (
        <footer>
            <div className="container">
                <div className="logo-container">
                    <a href="" className='logo'>
                        {svgIcons["logo"]}
                        <h4 className="registered-mark">© 2001—2024 Offliner</h4>
                    </a>
                </div>

                <Navigation />
            </div>
        </footer>
    );
}
