import React, { useState, useEffect } from 'react';
import "./Sales.css";
import gpu from '../../gpu';
import { addToWishList, removeFromWishList, getWishList } from '../../VishUtils';
import svgIcons from "../../svgIcons"

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
    let name = currentSaleItem.name;
    let imgUrl = currentSaleItem.imageURL;
    let specs = currentSaleItem.specs;
    let price = currentSaleItem.currentPrice;
    let origPrice = currentSaleItem.originalPrice;
    let discount = currentSaleItem.discount;

    const [isWished, setIsWished] = useState(false);

    useEffect(() => {
        const wishList = getWishList();
        setIsWished(wishList.some(item => item.name === name));
    }, [name]);

    const handleWishToggle = () => {
        if (isWished) {
            removeFromWishList(name);
        } else {
            addToWishList({ name, imgUrl, specs, price, origPrice, discount });
        }
        setIsWished(!isWished);
    }

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

            <div className="cost-container">
                <h2 className="cost">{currentSaleItem.currentPrice}</h2>
                <h4 className="previous-cost">{currentSaleItem.originalPrice || '—'}</h4>
            </div>

            <button className="add-to-vish" onClick={handleWishToggle}>
                {isWished ? svgIcons["addToWish"] : svgIcons["addToWishFilled"]}
            </button>

            <div className="sale-percent">{currentSaleItem.discount}</div>
        </section>
    );
}
