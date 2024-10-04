import React, { useState, useEffect, useRef } from 'react';
import "./Filter.css";
import svgIcons from "../../svgIcons";

export default function Filter() {
    const filter = useRef(null);
    function handleFilterClick(e) {
        e.preventDefault();    
        console.log('You clicked filter.');
    }

    return (
        <div ref={filter} className="filter" onClick={handleFilterClick}>
            Производитель
            {svgIcons["filterArrow"]}
        </div>
    );
}
