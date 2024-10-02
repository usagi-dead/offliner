import React from 'react';
import '../Navigation.css';
import svgIcons from '../../../svgIcons';

export default function SocialMedia({ href, icon, text }) {
    return (
        <a href={href} target="_blank" rel="noopener noreferrer" className='footer-button social-media-container'>
            {svgIcons[icon]}
            <h4 className="footer-text">{text}</h4>
        </a>
    );
}
