import React, { useState, useEffect } from 'react';
import { addToWishList, removeFromWishList, getWishList } from '../../VishUtils';
import svgIcons from "../../svgIcons";
import "./ToVish.css";

export default function ToVish({ vishItem }) {
    const { name, imgUrl, specs, price, origPrice, discount } = vishItem;

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
        setIsWished(prevState => !prevState);
    };

    return (
        <button 
            className="add-to-vish" 
            onClick={handleWishToggle} 
        >
            {isWished ? svgIcons["addToWishFilled"] : svgIcons["addToWish"]}
        </button>
    );
}
