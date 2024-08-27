import React, { useRef } from 'react';
import classes from './Card.module.css';
import Specs from './Specs/Specs';
import ToVish from "../Items/ToVish/ToVish";
import ToBusket from "../Items/ToBusket/ToBusket";

export default function Card({ product }) {
    const { name, imgUrl, specs, price, origPrice, discount, productUrl } = product;
    const cardRef = useRef(null);

    const handleCardClick = (event) => {
        if (!event.target.closest('button')) {
            window.location.href = `/gpu/product/${productUrl}`;
        }
    };

    return (
        <div className={classes.card} ref={cardRef} onClick={handleCardClick}>
            <div className={classes.cardImage} style={{ backgroundImage: `url(${imgUrl})` }}></div>
            <div className={classes.contentContainer}>
                <h1 className={classes.title}>{name}</h1>

                {origPrice == null ? (
                    <></>
                ) : (
                    <div className={classes.discount}>{discount}</div>
                )}

                <Specs specs={specs} />

                <div className={classes.bottomContainer}>
                    {price != null ? (
                        <>
                            <div className={classes.priceContainer}>
                                {origPrice == null ? (
                                    <h1 className={classes.price}>{price}</h1>
                                ) : (
                                    <>
                                        <h1 className={`${classes.price} ${classes.blue}`}>{price}</h1>
                                        <h3 className={classes.defaultPrice}>{origPrice}</h3>
                                    </>
                                )}
                            </div>
                        </>
                    ) : (
                        <h1 className={classes.notInStock}>Нет в наличии</h1>
                    )}
                    <div className={classes.buttons}>
                        <ToBusket />
                        <ToVish vishItem={product} />
                    </div>
                </div>
            </div>
        </div>
    );
}
