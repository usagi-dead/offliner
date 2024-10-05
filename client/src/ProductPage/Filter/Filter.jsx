import React, { useState, useEffect, useRef } from 'react';
import "./Filter.css";
import svgIcons from "../../svgIcons";

export default function Filter({ filterKey, text, isOpen }) {
    const filter = useRef(null);
    const [isClicked, setIsClicked] = useState(isOpen);
    const filtersss = ["zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc"];
    const isInput = true;
    const maxLength = "10";

    function handleFilterClick(e) {
        e.preventDefault();   
        setIsClicked(!isClicked);
    }

    function handleFilters(e) {
        e.stopPropagation();
    }

    return (
        <button ref={filter} className={`filter ${isClicked ? isInput ? "clicked-input" : "clicked" : ""}`} onClick={handleFilterClick}>
            <div className='filter-text'>
                {text}
                <span className='filter-arrow'>{svgIcons["smallArrow"]}</span>
            </div>
            {isInput ?
            <div className={`filters-input-container ${isClicked ? '' : 'filters-hidden'}`} onClick={handleFilters} tabIndex="-1">
                <input type="number" className="filter-input" placeholder='от' />
                <span className="filter-line" />
                <input type="number" className="filter-input" placeholder='8' />
            </div> :
            <div className={`filters ${isClicked ? '' : 'filters-hidden'}`} onClick={handleFilters} tabIndex="-1">
                {filtersss.map((prod, index) => {
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
