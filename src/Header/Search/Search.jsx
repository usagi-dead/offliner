import React, { useState } from 'react';
import "./Search.css"
import gpu from "../../gpu";
import InputSearchBar from './InputSearchBar/InputSearchBar';
import SearchButton from './SearchButton/SearchButton';
import SearchResults from './SearchResults/SearchResults';

export default function Search() {
    const [searchTerm, setSearchTerm] = useState('');
    const [searchResults, setSearchResults] = useState([]);

    const handleSearchChange = (event) => {
        setSearchTerm(event.target.value);
    };

    const handleSearch = () => {
        if (searchTerm) {
            const results = gpu.filter((item) =>
                item.name.toLowerCase().includes(searchTerm.toLowerCase())
            );
            setSearchResults(results);
        } else {
            setSearchResults([]);
        }
    };

    const handleKeyDown = (event) => {
        if (event.key === 'Enter') {
            handleSearch();
        } else if (event.key === 'Escape') {
            handleSearchClose();
        }
    };

    const handleSearchClose = () => {
        setSearchResults([]);
        setSearchTerm('');
    };

    return (
        <>
            <div className="search">
                <InputSearchBar 
                    searchTerm={searchTerm} 
                    onChange={handleSearchChange} 
                    onKeyDown={handleKeyDown} 
                />
                <SearchButton onClick={handleSearch} />
                {searchResults.length > 0 && (
                    <SearchResults results={searchResults} onClose={handleSearchClose} />
                )}
            </div>

            {searchResults.length > 0 && (
                <div className="overlay"></div>
            )}
        </>
    );
}
