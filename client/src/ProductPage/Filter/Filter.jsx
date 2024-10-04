import React, { useState, useEffect, useRef } from 'react';
import "./Filter.css";
import svgIcons from "../../svgIcons";

export default function Filter() {
    const filter = useRef(null);
    const [isClicked, setIsClicked] = useState(false);
    const filtersss = ["zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc"];

    function handleFilterClick(e) {
        e.preventDefault();   
        setIsClicked(!isClicked);
    }

    function handleFilterCheckboxClick(e) {
        e.stopPropagation();
    }

    return (
        <button ref={filter} className={isClicked ? "filter clicked" : "filter"} onClick={handleFilterClick}>
            <div className='filter-text'>
                Производитель
                <span className='filter-arrow'>{svgIcons["smallArrow"]}</span>
            </div>
            <div className={`filters ${isClicked ? '' : 'filters-hidden'}`} onClick={handleFilterCheckboxClick} tabIndex={isClicked ? "0" : "-1"}>
                {filtersss.map((prod, index) => {
                    return ( 
                        <div className='filter-checkbox-container'>
                            <input key={index} id={index} type='checkbox' className='filter-checkbox' tabIndex={isClicked ? "0" : "-1"} />
                            <label for={index}>{prod}</label>
                        </div>
                    );
                })}
            </div>
        </button>
    );
}
