import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './index.css';
import Header from './Header/Header';
import HomePage from './HomePage/HomePage';
import ProductPage from './ProductPage/ProductPage';
import VishPage from './VishPage/VishPage';
import Footer from './Footer/Footer';
import Product from './Product/Product';
import productNames from "./data";

export default function App() {
    return (
        <Router>
            <Header />
            <Routes>
                <Route path="/" element={<HomePage />} />
                <Route path="/favorites/:page" element={<VishPage />} />
                {productNames.map((product, index) => (
                    <Route 
                        key={index} 
                        path={`/${product.url}/:page?`} 
                        element={<ProductPage category={product.name} url={product.url} />} 
                    />
                ))}
                <Route path="/favorites" element={<VishPage />} />
                <Route 
                    path={`/:product/product/:productID`} 
                    element={<Product />} 
                />
            </Routes>
            <Footer />
        </Router>
    );
}
