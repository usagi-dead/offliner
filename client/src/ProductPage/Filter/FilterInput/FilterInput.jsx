import React from 'react';
import "./FilterInput.css";

export default function FilterInput({ max, isClicked }) {
    function handleFilters(e) {
        e.stopPropagation();
    }
    
    return (
        <div className={`filters-input-container ${isClicked ? 'filters-input-show' : 'filters-hidden'}`} onClick={handleFilters} tabIndex="-1">
            <input type="text" className="filter-input" placeholder='от' />
            <span className="filter-line" />
            <input type="text" className="filter-input" placeholder={max} />
        </div>
    );
}