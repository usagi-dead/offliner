import React, { useRef, useState, useEffect } from 'react';
import classes from './Card.module.css';
import { addToWishList, removeFromWishList, getWishList } from '../../VishUtils';
import { Link } from 'react-router-dom';

export default function Card({ name, imgUrl, specs, price, origPrice, discount, productUrl }) {
    const sections = {};
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

    Object.keys(specs).forEach((key) => {
        const [section, ...rest] = key.split(': ');
        if (!sections[section]) {
            sections[section] = [];
        }
        sections[section].push({ key: rest.join(': '), value: specs[key] });
    });

    const cardRef = useRef(null);

    return (
        <Link className={classes.card} ref={cardRef} to={`/gpu/product/${productUrl}`}>
            <div className={classes.cardImage} style={{ backgroundImage: `url(${imgUrl})` }}></div>
            <h1 className={classes.title}>{name}</h1>

            {origPrice == null ? (
                <></>
            ) : (
                <div className={classes.discount}>{discount}</div>
            )}

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
                        <button className={classes.toBusketButton}>В корзину</button>
                    </>
                ) : (
                    <h1 className={classes.notInStock}>Нет в наличии</h1>
                )}
            </div>

            <button className={classes.addToVish} onClick={handleWishToggle}>
                <svg width="56" height="56" viewBox="0 0 56 56" fill="none" xmlns="http://www.w3.org/2000/svg">
                <g clip-path="url(#clip0_11_36)">
                <mask id="path-1-inside-1_11_36" fill="white" >
                <path fill-rule="evenodd" clip-rule="evenodd" d="M5.6619 8.32391C11.8304 2.15546 21.8314 2.15545 27.9998 8.32388C27.9999 8.3239 27.9999 8.3239 27.9999 8.32388C27.9999 8.32386 28 8.32386 28 8.32388L28 8.32391C28 8.32394 28.0001 8.32394 28.0001 8.32391C34.1686 2.15545 44.1696 2.15545 50.3381 8.32391C56.5066 14.4924 56.5066 24.4934 50.3381 30.6619C50.3202 30.6798 50.3023 30.6976 50.2843 30.7154L32.2426 48.7572C29.8994 51.1003 26.1004 51.1003 23.7573 48.7572L5.69989 30.6998L5.6619 30.6619C-0.506569 24.4934 -0.506567 14.4924 5.6619 8.32391Z"/>
                </mask>
                {isWished ? <path fill-rule="evenodd" clip-rule="evenodd" d="M5.66196 8.3241C11.8304 2.15565 21.8314 2.15564 27.9999 8.32407C27.9999 8.32409 28 8.32409 28 8.32407V8.32407C28 8.32405 28 8.32405 28 8.32407L28.0001 8.3241C28.0001 8.32413 28.0002 8.32413 28.0002 8.3241V8.3241C34.1686 2.15564 44.1697 2.15564 50.3382 8.3241C56.5066 14.4926 56.5066 24.4936 50.3382 30.6621C50.3203 30.68 50.3024 30.6978 50.2844 30.7156L32.2426 48.7574C29.8995 51.1005 26.1005 51.1005 23.7574 48.7574L5.69995 30.6999L5.66196 30.6621C-0.506505 24.4936 -0.506503 14.4926 5.66196 8.3241Z" fill="var(--primary-color)" style={{transition: 'fill 0.3s'}}/> : <></>}
                <path d="M50.2843 30.7154L48.8768 29.2945L48.8701 29.3012L50.2843 30.7154ZM5.69989 30.6998L7.1141 29.2855L7.11175 29.2832L5.69989 30.6998ZM5.6619 30.6619L4.24768 32.0761L4.25003 32.0785L5.6619 30.6619ZM28 8.32391L26.5858 9.73813L28 8.32391ZM28 8.32388L26.5858 9.73809L28 8.32388ZM27.9998 8.32388L26.5856 9.7381L27.9998 8.32388ZM29.4141 6.90966C22.4645 -0.0398132 11.1972 -0.0397997 4.24768 6.9097L7.07611 9.73813C12.4635 4.35072 21.1982 4.35071 26.5856 9.7381L29.4141 6.90966ZM29.4142 6.9097L29.4142 6.90966L26.5858 9.73809L26.5858 9.73813L29.4142 6.9097ZM29.4143 9.73813C34.8017 4.35071 43.5365 4.35071 48.9239 9.73813L51.7523 6.9097C44.8028 -0.0398119 33.5354 -0.0398132 26.5859 6.9097L29.4143 9.73813ZM48.9239 9.73813C54.3113 15.1255 54.3113 23.8603 48.9239 29.2477L51.7523 32.0761C58.7018 25.1266 58.7018 13.8592 51.7523 6.9097L48.9239 9.73813ZM48.9239 29.2477C48.9082 29.2634 48.8925 29.279 48.8768 29.2946L51.6919 32.1362C51.7121 32.1163 51.7322 32.0962 51.7523 32.0761L48.9239 29.2477ZM33.6568 50.1714L51.6986 32.1296L48.8701 29.3012L30.8284 47.343L33.6568 50.1714ZM22.3431 50.1714C25.4673 53.2956 30.5326 53.2956 33.6568 50.1714L30.8284 47.343C29.2663 48.9051 26.7336 48.9051 25.1715 47.343L22.3431 50.1714ZM4.28567 32.114L22.3431 50.1714L25.1715 47.343L7.1141 29.2855L4.28567 32.114ZM4.25003 32.0785L4.28802 32.1163L7.11175 29.2832L7.07376 29.2453L4.25003 32.0785ZM4.24768 6.9097C-2.70183 13.8592 -2.70183 25.1266 4.24768 32.0761L7.07611 29.2477C1.6887 23.8603 1.6887 15.1255 7.07611 9.73813L4.24768 6.9097ZM26.5858 9.73813C27.3669 10.5192 28.6333 10.5192 29.4143 9.73813L26.5859 6.9097C27.3669 6.12868 28.6332 6.12868 29.4142 6.9097L26.5858 9.73813ZM29.529 9.61307C28.7738 10.5087 27.4141 10.5664 26.5858 9.73809L29.4142 6.90966C28.5858 6.08128 27.226 6.13904 26.4709 7.03469L29.529 9.61307ZM26.5856 9.7381C27.4019 10.5544 28.7628 10.5217 29.529 9.61307L26.4709 7.03469C27.237 6.12606 28.5978 6.09341 29.4141 6.90966L26.5856 9.7381Z" fill="var(--primary-color)" mask="url(#path-1-inside-1_11_36)" style={{transition: 'fill 0.3s'}}/>
                </g>
                <defs>
                <clipPath id="clip0_11_36">
                <rect width="56" height="56" fill="white"/>
                </clipPath>
                </defs>
                </svg>
            </button>
        </Link>
    );
}
