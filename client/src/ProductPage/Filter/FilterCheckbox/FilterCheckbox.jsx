import React from 'react';
import "./FilterCheckbox.css";

export default function FilterCheckbox({ filterKey, filters, isClicked }) {
    function handleFilters(e) {
        e.stopPropagation();
    }
    
    return (
        <div className={`filters ${isClicked ? 'filters-show' : 'filters-hidden'}`} onClick={handleFilters} tabIndex="-1">
            {filters.map((prod, index) => {
                return (
                    <label for={filterKey + "|" + index} className='filter-label'>
                        <input key={index} id={filterKey + "|" + index} type='checkbox' className='real-checkbox' tabIndex={isClicked ? "0" : "-1"} />
                        <span className='filter-checkbox' />
                        {prod}
                    </label>
                );
            })}
        </div>
    );
}
