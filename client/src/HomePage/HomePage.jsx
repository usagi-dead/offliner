import React, { useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import "./HomePage.css"
import Weather from './Weather/Weather'
import Sales from './Sales/Sales'
import Card from './Card/Card'
import productNames from "../data"
import svgIcons from "../svgIcons"

export default function HomePage() {
    const { pathname } = useLocation();

    useEffect(() => {
        window.scrollTo({ top: 0 });
    }, [pathname]);

    return (
        <>
            <section className="content-top">
                <Weather />
                <Sales />
            </section>

            <section className="content-products">
                <div className="title-container">
                    <h1 className="title">Наш ассортимент</h1>
                    {svgIcons["bottomArrow"]}
                </div>
            </section>

            <section className="cards-container">
                {productNames.map((product, index) => (
                <Card key={index} title={product.name } cardImage={product.file} urlProduct={product.url} />
                ))}
            </section>
        </>
    )
}