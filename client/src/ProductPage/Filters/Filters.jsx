import React, { useRef } from 'react';
import "./Filters.css";
import Filter from '../Filter/Filter';
import FilterButton from '../FilterButton/FilterButton';

export default function Filters() {
    const filters = useRef(null);

    return (
        <div>
            <div ref={filters} className='filters-container'>
                <Filter />
                <Filter />
                <Filter />
                <Filter />
                <Filter />
                <Filter />
                <Filter />
                <Filter />
                <Filter />
                <Filter />
                <Filter />
                <Filter />
                <Filter />
                <Filter />
            </div>
            <FilterButton filters={filters} />
        </div>
    );
}