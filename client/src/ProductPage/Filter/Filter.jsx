import React, { useState, useEffect, useRef } from 'react';
import "./Filter.css";
import svgIcons from "../../svgIcons";
import FilterCheckbox from "./FilterCheckbox/FilterCheckbox"
import FilterInput from "./FilterInput/FilterInput"

export default function Filter({ filterKey, text, isOpen, isInput, filters, max }) {
    const filter = useRef(null);
    const [isClicked, setIsClicked] = useState(isOpen);

    function handleFilterClick(e) {
        e.preventDefault();
        setIsClicked(!isClicked);
    }

    return (
        <button ref={filter} className={`filter ${isClicked ? "clicked" : ""}`} onClick={handleFilterClick}>
            <div className='filter-title-container'>
                <span className='filter-text'>{text}</span>
                <span className='filter-arrow'>{svgIcons["smallArrow"]}</span>
            </div>
            {isInput ? <FilterInput max={max} isClicked={isClicked} /> : <FilterCheckbox filterKey={filterKey} filters={filters} isClicked={isClicked} />}
        </button>
    );
}
