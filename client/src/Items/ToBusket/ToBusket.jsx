import React from 'react';
import "./ToBusket.css"

export default function ToBusket() {
    const handleAddToBusket = (event) => {
        event.stopPropagation();
        event.preventDefault(); 
    };

    return (
        <button 
            className="add-to-busket" 
            onClick={handleAddToBusket}
        >
            В корзину
        </button>
    );
}
