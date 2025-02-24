import './SearchResultsCard.css';
import React from 'react';
import { Link } from 'react-router-dom';
import Specs from "../../../../Card/Specs/Specs";

export default function SearchResults({ item, index, isLast }) {
    return (
        <div className='relative'>
            <li key={index} className="search-result-item">
                <Link to={`/gpu/${item.name}`} className='search-result-text'>
                    <div className='search-img-container'>
                        <img src={item.imageURL} alt={item.name} />
                    </div>
                    <div className="search-container">
                        <span className='name'>{item.name}</span>

                        <div className='search-specs-container'>
                            <Specs textSize="16" specs={item.specs} lineNum={1} />
                        </div>

                        <div className='cost-container'>
                            {item.originalPrice == null ?
                                <h1 className="search-price">{item.currentPrice}</h1> :
                                <>
                                    <h1 className={"search-price blue"}>{item.currentPrice}</h1>
                                    <h3 className="search-default-price">{item.originalPrice}</h3>
                                    <div className="search-discount">{item.discount}</div>
                                </>
                            }
                        </div>
                    </div>
                </Link>
            </li>
            {!isLast && <div className='search-line'></div>}
        </div>
    );
}
