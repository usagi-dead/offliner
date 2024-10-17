import React, { useState } from 'react';
import "./FilterInput.css";

export default function FilterInput({ max, isClicked }) {
    const [isActive, setIsActive] = useState(false); 

    function handleFilters(e) {
        e.stopPropagation();
    }

    function handleFocus() {
        setIsActive(true); 
    }

    function handleBlur() {
        setIsActive(false); 
    }

    return (
        <div 
            className={`filters-input-container ${isClicked ? 'filters-input-show' : 'filters-hidden'} ${isActive ? 'filters-input-active' : ''}`} 
            onClick={handleFilters} 
            tabIndex="-1"
        >
            <input 
                type="text" 
                className="filter-input" 
                placeholder='от' 
                onFocus={handleFocus} 
                onBlur={handleBlur}
            />
            <span className="filter-line" />
            <input 
                type="text" 
                className="filter-input" 
                placeholder={max} 
                onFocus={handleFocus} 
                onBlur={handleBlur}
            />
        </div>
    );
}
