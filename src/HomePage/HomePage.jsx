import "./HomePage.css"
import Weather from './Weather/Weather'
import Sales from './Sales/Sales'
import Card from './Card/Card'
import productNames from "../data"

export default function HomePage() {
    return (
        <>
            <section className="content-top">
                <Weather />
                <Sales />
            </section>

            <section className="content-products">
                <div className="title-container">
                    <h1>Наш ассортимент</h1>

                    <svg width="58" height="71" viewBox="0 0 58 71" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M1 42L21.8666 63.2393C25.7857 67.2283 32.2143 67.2283 36.1334 63.2393L57 42" stroke="#025ADD" stroke-width="2" stroke-linecap="round"/>
                    <path d="M29 66L29 1" stroke="#025ADD" stroke-width="2" stroke-linecap="round"/>
                    </svg>
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