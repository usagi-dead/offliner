import React from 'react';
import './Navigation.css';
import Contact from './Buttons/Contact';
import SocialMedia from './Buttons/SocialMedia';

export default function Navigation() {
    return (
        <div className="footer-content-container">
            <h2 className="footer-text footer-title">Контакты:</h2>

            <Contact text="offliner@gmail.com" />
            <Contact text="+375 (33) 123-45-67" />

            <h2 className="footer-text footer-title">Социальные сети:</h2>

            <SocialMedia href="https://www.instagram.com/" icon="inst" text="Instagram" />
            <SocialMedia href="https://www.youtube.com/" icon="yt" text="YouTube" />
            <SocialMedia href="https://x.com/?lang=ru" icon="x" text="Twitter" />
            <SocialMedia href="https://web.telegram.org/k/" icon="tg" text="Telegram" />
        </div>
    );
}
