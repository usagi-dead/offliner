import React, { useState, useEffect } from 'react';
import "./FilterButton.css";

export default function FilterButton({ filters }) {
    const [isFixed, setIsFixed] = useState(false);

    const handleScroll = () => {
        const filtersBottom = filters.current.getBoundingClientRect().bottom;

        if (filtersBottom > window.scrollY) {
            setIsFixed(true);
        } else {
            setIsFixed(false);
        }
    };

    useEffect(() => {
        window.addEventListener('scroll', handleScroll);
    
        return () => {
          window.removeEventListener('scroll', handleScroll);
        };
      }, []);

    return (
        <div className={isFixed ? "filter-button fixed" : "filter-button"}>
            Найти
        </div>
    );
}
