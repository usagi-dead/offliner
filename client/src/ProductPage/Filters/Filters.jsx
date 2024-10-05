import React, { useRef } from 'react';
import "./Filters.css";
import Filter from '../Filter/Filter';
import FilterButton from '../FilterButton/FilterButton';

export default function Filters() {
    const filters = useRef(null);
    const filtersss = [
        {name: "Производитель", isOpen: true},
        {name: "Производитель", isOpen: true},
        {name: "Производитель", isOpen: false},
        {name: "Производитель", isOpen: false},
        {name: "Производитель", isOpen: false},
        {name: "Производитель", isOpen: false},
        {name: "Производитель", isOpen: false},
        {name: "Производитель", isOpen: false},
        {name: "Производитель", isOpen: false},
        {name: "Производитель", isOpen: false},
        {name: "Производитель", isOpen: false},
        {name: "Производитель", isOpen: false}
    ];

    return (
        <div>
            <div ref={filters} className='filters-container'>
                {filtersss.map((prod, index) => {
                    return ( 
                        <Filter filterKey={index} text={prod.name} isOpen={prod.isOpen} />
                    );
                })}
            </div>
            <FilterButton filters={filters} />
        </div>
    );
}