import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import "./Navigation.css"
import productNames from "../../data";

export default function Navigation() {
    const location = useLocation();

    return (
        <div className="nav-container">
            <ul className="nav-items">
                {productNames.map((product, index) => {
                    const isActive = location.pathname.startsWith(`/${product.url}`);
                    return (
                        <li key={index} className={isActive ? 'nav-item active' : 'nav-item'}>
                            <Link to={`/${product.url}/1`}>{product.name}</Link>
                        </li>
                    );
                })}
            </ul>
        </div>
    )
}
