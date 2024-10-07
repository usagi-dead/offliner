import React, { useRef } from 'react';
import "./Filters.css";
import Filter from '../Filter/Filter';
import FilterButton from '../FilterButton/FilterButton';

export default function Filters() {
    const filters = useRef(null);
    const filtersss = [
        {
            name: "Производитель", 
            isOpen: true,
            isInput: false,
            filters: ["zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc"],
            max: null,
        },
        {
            name: "Цена ввыааааааа ываываываыв аываываываыв", 
            isOpen: true,
            isInput: true,
            filters: null,
            max: 5,
        },
        {
            name: "Производитель", 
            isOpen: false,
            isInput: false,
            filters: ["zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc"],
            max: null,
        },
        {
            name: "Производитель", 
            isOpen: false,
            isInput: false,
            filters: ["zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc"],
            max: null,
        },
        {
            name: "Производитель", 
            isOpen: false,
            isInput: false,
            filters: ["zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc"],
            max: null,
        },
        {
            name: "Производитель", 
            isOpen: false,
            isInput: false,
            filters: ["zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc"],
            max: null,
        },
        {
            name: "Производитель", 
            isOpen: false,
            isInput: false,
            filters: ["zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc","zxc"],
            max: null,
        },
    ];

    return (
        <div>
            <div ref={filters} className='filters-container'>
                {filtersss.map((filter, index) => {
                    return ( 
                        <Filter 
                            filterKey={index} 
                            text={filter.name} 
                            isOpen={filter.isOpen} 
                            isInput={filter.isInput}
                            filters={filter.filters}
                            max={filter.max}
                        />
                    );
                })}
            </div>
            <FilterButton filters={filters} />
        </div>
    );
}