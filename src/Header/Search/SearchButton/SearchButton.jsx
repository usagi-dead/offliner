import './SearchButton.css';
import React from 'react';
import svgIcons from '../../../svgIcons';

export default function SearchButton({ onClick }) {
    return (
        <div className="search-button" onClick={onClick}>
            {svgIcons["search"]}
        </div>
    );
}
