import React, { useState } from 'react';
import './Navigation.css';
import svgIcons from '../../svgIcons';

export default function Navigation() {
    const [copiedText, setCopiedText] = useState('');

    const copyToClipboard = (text) => {
        navigator.clipboard.writeText(text)
            .then(() => {
                setCopiedText(text);
                setTimeout(() => {
                    setCopiedText('');
                }, 3000);
            })
            .catch((err) => console.error('Failed to copy text: ', err));
    };

    return (
        <>
            <div className="footer-content-container">
                <h2 className="footer-text footer-title">
                    Контакты:
                </h2>

                <div className="footer-button" onClick={() => copyToClipboard('offliner@gmail.com')}>
                    <h4 className="footer-text">offliner@gmail.com</h4>
                </div>

                <div className="footer-button" onClick={() => copyToClipboard('+375 (33) 123-45-67')}>
                    <h4 className="footer-text">+375 (33) 123-45-67</h4>
                </div>

                <h2 className="footer-text footer-title">
                    Социальные сети:
                </h2>

                <a href="https://www.instagram.com/" target="_blank" rel="noopener noreferrer" className='footer-button social-media-container'>
                    {svgIcons["inst"]}
                    <h4 className="footer-text">Instagram</h4>
                </a>

                <a href="https://www.youtube.com/" target="_blank" rel="noopener noreferrer" className='footer-button social-media-container'>
                    {svgIcons["yt"]}
                    <h4 className="footer-text">YouTube</h4>
                </a>

                <a href="https://x.com/?lang=ru" target="_blank" rel="noopener noreferrer" className='footer-button social-media-container'>
                    {svgIcons["x"]}
                    <h4 className="footer-text">Twitter</h4>
                </a>

                <a href="https://web.telegram.org/k/" target="_blank" rel="noopener noreferrer" className='footer-button social-media-container'>
                    {svgIcons["tg"]}
                    <h4 className="footer-text">Telegram</h4>
                </a>
            </div>
            {copiedText && (
                <div className="notification">
                    {`${copiedText} скопировано в буфер обмена`}
                </div>
            )}
        </>
    );
}
