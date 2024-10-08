import React, { useState, useEffect, useRef } from 'react';
import "./Filter.css";
import svgIcons from "../../svgIcons";

export default function Filter({ filterKey, text, isOpen, isInput, filters, max }) {
    const filter = useRef(null);
    const [isClicked, setIsClicked] = useState(isOpen);

    function handleFilterClick(e) {
        e.preventDefault();
        setIsClicked(!isClicked);
    }

    function handleFilters(e) {
        e.stopPropagation();
    }

    return (
        <button ref={filter} className={`filter ${isClicked ? "clicked" : ""}`} onClick={handleFilterClick}>
            <div className='filter-title-container'>
                <span className='filter-text'>{text}</span>
                <span className='filter-arrow'>{svgIcons["smallArrow"]}</span>
            </div>
            {isInput ?
            <div className={`filters-input-container ${isClicked ? 'filters-input-show' : 'filters-hidden'}`} onClick={handleFilters} tabIndex="-1">
                <input type="text" className="filter-input" placeholder='от' />
                <span className="filter-line" />
                <input type="text" className="filter-input" placeholder={max} />
            </div> :
            <div className={`filters ${isClicked ? 'filters-show' : 'filters-hidden'}`} onClick={handleFilters} tabIndex="-1">
                {filters.map((prod, index) => {
                    return ( 
                        <label for={filterKey+"|"+index} className='filter-label'>
                            <input key={index} id={filterKey+"|"+index} type='checkbox' className='real-checkbox' tabIndex={isClicked ? "0" : "-1"} />
                            <span className='filter-checkbox' />
                            {prod}
                        </label>
                    );
                })}
            </div>
            }
        </button>
    );
}
