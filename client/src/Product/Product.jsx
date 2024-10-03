import React, { useEffect } from 'react';
import "./Product.css";
import { useParams, useLocation } from 'react-router-dom';
import gpu from '../gpu';
import Specs from "../Card/Specs/Specs";
import ToVish from "../Items/ToVish/ToVish";
import ToBusket from "../Items/ToBusket/ToBusket";

export default function ProductPage() {
    const { productID } = useParams();
    const product = gpu[productID];
    const { pathname } = useLocation();

    if (!product) {
        return <div className="product-page product-container title">Товар не найден</div>;
    }

    const groupedSpecs = Object.entries(product.specs).reduce((acc, [key, value]) => {
        const [category, spec] = key.split(':').map(part => part.trim());
        if (!acc[category]) {
            acc[category] = [];
        }
        acc[category].push({ spec, value });
        return acc;
    }, {});

    useEffect(() => {
        window.scrollTo({ top: 0, behavior: 'smooth' });
    }, [pathname]);

    return (
        <div className="product-page">
            <div className='product-container'>
                <h1 className='title'>Видеокарты</h1>
                <div className='product-content-container'>
                    <div className='product-image-container'>
                        <div className='image-swapper'>
                            <img src={product.imageURL} alt={product.name} className='swap-image' />
                            <img src={product.imageURL} alt={product.name} className='swap-image' />
                            <img src={product.imageURL} alt={product.name} className='swap-image' />
                            <img src={product.imageURL} alt={product.name} className='swap-image' />
                            <img src={product.imageURL} alt={product.name} className='swap-image' />
                        </div>
                        <div className='main-image'>
                            <img src={product.imageURL} alt={product.name} className='product-image' />
                        </div>
                    </div>
                    <div className='product-text-container'>
                        <h2 className='product-name'>{product.name}</h2>

                        <Specs specs={product.specs} textSize="20" />

                        <div className="product-prices">
                            {product.currentPrice && <h2 className={product.originalPrice ? "price blue" : "price"}>{product.currentPrice}</h2>}
                            {product.originalPrice && 
                            <h3 className='original-price'>
                                {product.originalPrice}
                                {product.discount && <h4 className='discount'>{product.discount}</h4>}
                            </h3>}
                        </div>

                        <div className='product-buttons-container'>
                            <ToBusket />
                            <ToVish vishItem={product} />
                        </div>
                    </div>
                </div>
                <div className="product-specs">
                    <h1 className='title'>Характеристики</h1>
                    {Object.entries(groupedSpecs).map(([category, specs]) => (
                        <div key={category}>
                            <h2 className="specs-category">{category}</h2>
                            {specs.map((spec, index) => (
                                <div
                                    key={spec.spec}
                                    className={`spec-container ${index === specs.length - 1 ? 'last-spec' : ''}`}
                                >
                                    <div className='spec-name'>{spec.spec}:</div>
                                    <div className='spec-value'>{spec.value}</div>
                                </div>
                            ))}
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}
