import './InputSearchBar.css';
import React from 'react';

export default function InputSearchBar({ searchTerm, onChange, onKeyDown }) {
    return (
        <input
            type="text"
            className="search-bar"
            placeholder="Поиск"
            value={searchTerm}
            onChange={onChange}
            onKeyDown={onKeyDown}
        />
    );
}
