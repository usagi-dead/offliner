import React from 'react';
import { Link } from 'react-router-dom';
import classes from './Card.module.css';
import Specs from './Specs/Specs';
import ToVish from "../Items/ToVish/ToVish";
import ToBusket from "../Items/ToBusket/ToBusket";

export default function Card({ product }) {
    const { name, imgUrl, specs, price, origPrice, discount, productUrl } = product;

    return (
        <Link className={classes.card} to={`/gpu/product/${productUrl}`}>
            <div className={classes.cardImage} style={{ backgroundImage: `url(${imgUrl})` }}></div>
            <div className={classes.contentContainer}>
                <h1 className={classes.title}>{name}</h1>

                {origPrice == null ? (
                    <></>
                ) : (
                    <div className={classes.discount}>{discount}</div>
                )}

                <Specs specs={specs} textSize="16" lineNum={3} />

                <div className={classes.bottomContainer}>
                    {price != null ? (
                        <>
                            <div className={classes.priceContainer}>
                                {origPrice == null ? (
                                    <h1 className={classes.price}>{price}</h1>
                                ) : (
                                    <>
                                        <h3 className={classes.defaultPrice}>{origPrice}</h3>
                                        <h1 className={`${classes.price} ${classes.blue}`}>{price}</h1>
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
        </Link>
    );
}
