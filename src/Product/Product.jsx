import React from 'react';
import "./Product.css";
import { useParams } from 'react-router-dom';
import gpu from '../gpu';

export default function ProductPage() {
    const { productID } = useParams();
    const product = gpu[productID];

    if (!product) {
        return <div>Товар не найден</div>;
    }

    return (
        <div className="product-page">
            <div className='product-container'>
                <h1 className='title'>{product.name}</h1>
                <div className='product-content-container'>
                    <div className='product-image-container'>
                        <img src={product.imageURL} alt={product.name} className='product-image' />
                    </div>
                    <div className="product-prices">
                        {product.currentPrice && <h2 className='price'>{product.currentPrice}</h2>}
                        {product.originalPrice && <h3 className='original-price'>{product.originalPrice}</h3>}
                        {product.discount && <h4 className='discount'>{product.discount}</h4>}
                    </div>
                </div>
                <div className="product-specs">
                    {Object.entries(product.specs).map(([key, value]) => (
                        <div key={key}>
                            <strong>{key}:</strong> {value}
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}
