import React, { useState, useEffect, useRef } from 'react';
import "./Filter.css";
import svgIcons from "../../svgIcons";

export default function Filter() {
    const filter = useRef(null);
    const [isClicked, setIsClicked] = useState(false);

    function handleFilterClick(e) {
        e.preventDefault();   
        console.log(e);

        setIsClicked
    }

    return (
        <button ref={filter} className={isFixed ? "filter clicked" : "filter"} onClick={handleFilterClick}>
            Производитель
            {svgIcons["filterArrow"]}
        </button>
    );
}
