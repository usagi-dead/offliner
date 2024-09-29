import React, { useState, useEffect } from 'react';
import "./Header.css";
import HeaderButton from './HeaderButton/HeaderButton';
import Search from './Search/Search';
import Navigation from './Navigation/Navigation';
import { Link } from 'react-router-dom';
import svgIcons from '../svgIcons';

export default function Header() {
    const [isScrolled, setIsScrolled] = useState(false);

    useEffect(() => {
        const handleScroll = () => {
            const scrollTop = window.scrollY;
            setIsScrolled(scrollTop > 0); 
        };

        window.addEventListener("scroll", handleScroll);

        return () => {
        window.removeEventListener("scroll", handleScroll);
        };
    }, []);

    return (
        <header>
            <div className={isScrolled ? "container scrolled" : "container"}>
                <Link to={`/`} className='logo'>
                    {svgIcons["logo"]}
                </Link>

                <Search />

                <div className='buttons-container'>

                    <HeaderButton 
                        svg="theme"
                    />

                    <HeaderButton 
                        link="/"
                        svg="profile"
                    />

                    <HeaderButton 
                        link="/favorites/1"
                        svg="vish"
                    />

                    <HeaderButton 
                        link="/"
                        svg="busket"
                    />

                </div>

            </div>

            <Navigation />
            
        </header>
    );
}