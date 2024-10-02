import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import "./ProductPage.css";
import Card from '../Card/Card';
import gpu from '../gpu';
import Pagination from './Pagination/Pagination';

const ITEMS_PER_PAGE = 27;

export default function ProductPage({ url, category }) {
    const { page } = useParams();
    const navigate = useNavigate();
    const initialPage = parseInt(page) || 1;
    const [currentPage, setCurrentPage] = useState(initialPage);
    const items = gpu;
    const totalPages = Math.ceil(items.length / ITEMS_PER_PAGE);
    const { pathname } = useLocation();

    useEffect(() => {
        if (page !== currentPage.toString()) {
            navigate(`/${url.toLowerCase()}/${currentPage}`, { replace: true });
        }
    }, [currentPage, url, navigate, page]);

    const currentItems = items.slice((currentPage - 1) * ITEMS_PER_PAGE, currentPage * ITEMS_PER_PAGE);

    const handlePageChange = (page) => {
        setCurrentPage(page);
    };

    useEffect(() => {
        window.scrollTo({ top: 0, behavior: 'smooth' });
    }, [pathname]);


    const [isScrolled, setIsScrolled] = useState(false);

    useEffect(() => {
        const handleScroll = () => {
            const scrollTop = window.scrollY;
            setIsScrolled(scrollTop > 20); 
        };

        window.addEventListener("scroll", handleScroll);

        return () => {
        window.removeEventListener("scroll", handleScroll);
        };
    }, []);

    return (
        <section className="products-page">
            <div className='products-container'>
                <h1 className='title'>{category}</h1>
                <div className='products-content-container'>
                    <div className='products-filters-container'>

                    </div>
                    
                    <div className='products-cards-container'>
                        {currentItems.map((item, index) => (
                            <Card 
                                key={index} 
                                product={{
                                    name: item.name, 
                                    imgUrl: item.imageURL, 
                                    specs: item.specs, 
                                    price: item.currentPrice,
                                    origPrice: item.originalPrice,
                                    discount: item.discount,
                                    productUrl: index + (ITEMS_PER_PAGE * (page - 1))
                                }} 
                            />
                        ))}
                    </div>
                </div>
                {totalPages > 1 && (
                    <Pagination 
                        currentPage={currentPage} 
                        totalPages={totalPages} 
                        onPageChange={handlePageChange} 
                    />
                )}
            </div>
        </section>
    );
}
