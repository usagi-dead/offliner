import React, { useState } from 'react';
import '../Navigation.css';

export default function Contact({ text }) {
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
            <div className="footer-button" onClick={() => copyToClipboard(text)}>
                <h4 className="footer-text">{text}</h4>
            </div>
            {copiedText && (
                <div className="notification">
                    {`${copiedText} скопировано в буфер обмена`}
                </div>
            )}
        </>
    );
}
