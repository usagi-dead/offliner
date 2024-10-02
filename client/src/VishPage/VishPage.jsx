import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import "./VishPage.css";
import Card from '../Card/Card';
import { getWishList } from '../VishUtils';

export default function VishPage() {
    const { pathname } = useLocation();
    const [items, setItems] = useState([]);

    useEffect(() => {
        window.scrollTo({ top: 0 });
    }, [pathname]);

    useEffect(() => {
        setItems(getWishList());
    }, []);

    return (
        <section className="vish-page">
            <div className='vish-container'>
                <h1 className='title'>Желаемое</h1>
                <div className='vish-cards-container'>
                    {items.map((item, index) => (
                        <Card 
                            key={index} 
                            product={{
                                name: item.name, 
                                imgUrl: item.imgUrl, 
                                specs: item.specs, 
                                price: item.price,
                                origPrice: item.origPrice,
                                discount: item.discount,
                                productUrl: index
                            }} 
                        />
                    ))}
                </div>
            </div>
        </section>
    );
}
