import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import "./VishPage.css";
import Card from '../Card/Card';
import { getWishList } from '../VishUtils';

export default function VishPage() {
    const { pathname } = useLocation();
    const [items, setItems] = useState([]);

    useEffect(() => {
        setItems(getWishList());
    }, []);

    useEffect(() => {
        window.scrollTo({ top: 0, behavior: 'smooth' });
    }, [pathname]);

    return (
        <section className="vish-page">
            <div className='vish-container'>
                <h1 className='title'>Желаемое</h1>
                <div className='vish-cards-container'>
                    {items.map((card, index) => (
                        <Card 
                            key={index} 
                            name={card.name} 
                            imgUrl={card.imgUrl} 
                            specs={card.specs} 
                            price={card.price}  
                            origPrice={card.origPrice} 
                            discount={card.discount}
                        />
                    ))}
                </div>
            </div>
        </section>
    );
}
