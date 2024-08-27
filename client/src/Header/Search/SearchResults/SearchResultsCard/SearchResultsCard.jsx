import './SearchResultsCard.css';
import React from 'react';
import { Link } from 'react-router-dom';

export default function SearchResults({ item, index }) {
    return (
        <li key={index} className="search-result-item">
            <Link to={`/gpu/${item.name}`} className='search-result-text'>
                <div className='search-img-container'>
                    <img src={item.imageURL} alt={item.name} />
                </div>
                <div className="search-container">
                    <span className='name'>{item.name}</span>
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
            <div className='search-line'></div>
        </li>
    );
}
