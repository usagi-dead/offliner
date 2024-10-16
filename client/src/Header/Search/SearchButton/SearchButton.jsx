import './SearchButton.css';
import React from 'react';
import svgIcons from '../../../svgIcons';

export default function SearchButton({ onClick, background }) {
    return (
        <div className={"search-button" + background} onClick={onClick}>
            {svgIcons["search"]}
        </div>
    );
}
