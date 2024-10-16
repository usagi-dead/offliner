import './SearchResults.css';
import React from 'react';
import SearchResultsCard from './SearchResultsCard/SearchResultsCard';

export default function SearchResults({ results }) {
    return (
        <div className='search-items-container'>
            <div className="search-results">
                <div className='search-padding' />

                <ul>
                    {results.map((item, index) => (
                        <SearchResultsCard item={item} index={index} isLast={index === results.length - 1}  />
                    ))}
                </ul>

                <div className='search-padding' />
            </div>
        </div>
    );
}
