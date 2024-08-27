import React, { useState } from 'react';
import "./Sales.css";
import gpu from '../../gpu';
import svgIcons from "../../svgIcons";
import ToVish from "../../Items/ToVish/ToVish";

export default function Sales() {
    const discountedItems = gpu.filter(item => item.discount !== null);
    const [currentIndex, setCurrentIndex] = useState(0);
    const [isAnimating, setIsAnimating] = useState(false);
    const [animationDirection, setAnimationDirection] = useState('');

    const handleAnimationEnd = () => {
        setIsAnimating(false);
    };

    const handlePrevious = () => {
        if (isAnimating) return;
        setIsAnimating(true);
        setAnimationDirection('left'); 
        setCurrentIndex(prevIndex => (prevIndex === 0 ? discountedItems.length - 1 : prevIndex - 1));
    };

    const handleNext = () => {
        if (isAnimating) return;
        setIsAnimating(true);
        setAnimationDirection('right');
        setCurrentIndex(prevIndex => (prevIndex === discountedItems.length - 1 ? 0 : prevIndex + 1));
    };

    const currentSaleItem = discountedItems[currentIndex];

    return (
        <section className="sales-container">
            <h1 className="title">Акции</h1>
                
            <div className="sale-items-wrapper">
                <button className="arrow left-arrow" onClick={handlePrevious}>
                    {svgIcons["leftArrow"]}
                </button>

                <div className={`sale-items ${isAnimating ? (animationDirection === 'left' ? 'slide-left' : 'slide-right') : ''}`} onAnimationEnd={handleAnimationEnd}>
                    <div className="item">
                        <div className='sales-img-container'>
                            <img src={currentSaleItem.imageURL} alt="" className="sale-image" />
                        </div>
                        <h4 className="item-text">{currentSaleItem.name}</h4>
                    </div>
                </div>

                <button className="arrow right-arrow" onClick={handleNext}>
                    {svgIcons["rightArrow"]}
                </button>
            </div>

            <div className='sale-bot-container'>
                <div className="cost-container">
                    <h2 className="cost">{currentSaleItem.currentPrice}</h2>
                    <h4 className="previous-cost">{currentSaleItem.originalPrice || '—'}</h4>
                </div>

                <ToVish vishItem={currentSaleItem} />
            </div>

            <div className="sale-percent">{currentSaleItem.discount}</div>
        </section>
    );
}
