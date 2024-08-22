import React from 'react';
import "./Header.css";
import HeaderButton from './HeaderButton/HeaderButton';
import ThemeSwitchButton from './ThemeSwitchButton/ThemeSwitchButton';
import Search from './Search/Search';
import Navigation from './Navigation/Navigation';
import { Link } from 'react-router-dom';
import svgIcons from '../svgIcons';

export default function Header() {
    return (
        <header>
            <div className="container">
                <Link to={`/`}>
                    {svgIcons["logo"]}
                </Link>

                <Search />

                <div className='buttons-container'>

                    <ThemeSwitchButton />

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
                        isBusket="1"
                    />

                </div>

            </div>

            <Navigation />
            
        </header>
    );
}