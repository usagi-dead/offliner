import './SearchResults.css';
import React from 'react';
import svgIcons from '../../../svgIcons';
import { Link } from 'react-router-dom';
import SearchResultsCard from './SearchResultsCard/SearchResultsCard';

export default function SearchResults({ results, onClose }) {
    return (
        <div className='search-items-container'>
            <div className='title-search-container'>
                <h1 className='search-title'>Результаты поиска</h1>
                <button className='close-button-search' onClick={onClose}>
                    {svgIcons["close"]}
                </button>
            </div>
            <div className="search-results">
                <ul>
                    {results.map((item, index) => (
                        <SearchResultsCard item={item} index={index} />
                    ))}
                </ul>
            </div>
        </div>
    );
}
