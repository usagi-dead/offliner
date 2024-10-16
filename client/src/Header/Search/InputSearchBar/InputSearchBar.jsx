import './InputSearchBar.css';
import React from 'react';

export default function InputSearchBar({ searchTerm, onChange, onKeyDown, background }) {
    return (
        <input
            type="text"
            className={"search-bar" + background}
            placeholder="Поиск"
            value={searchTerm}
            onChange={onChange}
            onKeyDown={onKeyDown}
        />
    );
}
