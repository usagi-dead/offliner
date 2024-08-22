import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import "./ProductPage.css";
import Card from './Card/Card';
import gpu from '../gpu';

const ITEMS_PER_PAGE = 27;

export default function ProductPage({ url, category }) {
    const { page } = useParams();
    const navigate = useNavigate();
    const initialPage = parseInt(page) || 1;
    const [currentPage, setCurrentPage] = useState(initialPage);
    const items = gpu;
    const totalPages = Math.ceil(items.length / ITEMS_PER_PAGE);
    const { pathname } = useLocation();
    const maxPagesToShow = 9;

    useEffect(() => {
        if (page !== currentPage.toString()) {
            navigate(`/${url.toLowerCase()}/${currentPage}`, { replace: true });
        }
    }, [currentPage, url, navigate, page]);

    const currentItems = items.slice((currentPage - 1) * ITEMS_PER_PAGE, currentPage * ITEMS_PER_PAGE);

    const handlePageChange = (page) => {
        if (page > 0 && page <= totalPages)
            setCurrentPage(page);
    };

    useEffect(() => {
        window.scrollTo({ top: 0, behavior: 'smooth' });
    }, [pathname]);

    const renderPaginationButtons = () => {
        const buttons = [];
        let startPage, endPage;

        if (totalPages <= maxPagesToShow) {
            startPage = 1;
            endPage = totalPages;
        } else {
            startPage = Math.max(currentPage - 4, 1);
            endPage = Math.min(currentPage + 4, totalPages);

            if (startPage === 1) {
                endPage = maxPagesToShow;
                startPage = 1;
            } 

            if (endPage === totalPages) {
                startPage = totalPages - maxPagesToShow + 1;
                endPage = totalPages;
            }
        }

        if (startPage > 1) {
            buttons.push(
                <button
                    key={1}
                    className={currentPage === 1 ? 'pagination-button pagination-active' : 'pagination-button'}
                    onClick={() => handlePageChange(1)}
                >
                    1
                </button>
            );
            if (startPage > 2) {
                buttons.push(<span key="start-ellipsis" className='dots'>...</span>);
            }
        }

        for (let i = startPage; i <= endPage; i++) {
            buttons.push(
                <button
                    key={i}
                    className={currentPage === i ? 'pagination-button pagination-active' : 'pagination-button'}
                    onClick={() => handlePageChange(i)}
                >
                    {i}
                </button>
            );
        }

        if (endPage < totalPages) {
            if (endPage < totalPages - 1) {
                buttons.push(<span key="end-ellipsis" className='dots'>...</span>);
            }
            buttons.push(
                <button
                    key={totalPages}
                    className={currentPage === totalPages ? 'pagination-button pagination-active' : 'pagination-button'}
                    onClick={() => handlePageChange(totalPages)}
                >
                    {totalPages}
                </button>
            );
        }

        return buttons;
    };

    return (
        <section className="products-page">
            <div className='products-container'>
                <h1 className='title'>{category}</h1>
                <div className='products-cards-container'>
                    {currentItems.map((card, index) => (
                        <Card 
                            key={index} 
                            name={card.name} 
                            imgUrl={card.imageURL} 
                            specs={card.specs} 
                            price={card.currentPrice}
                            origPrice={card.originalPrice}
                            discount={card.discount}
                            productUrl={index}
                        />
                    ))}
                </div>
                {totalPages > 1 && (
                    <div className="pagination-container">
                        <button className='move-button' onClick={() => handlePageChange(currentPage - 1)}>
                            <svg width="7" height="14" viewBox="0 0 7 14" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M6 0.5L1.25221 5.64344C0.545018 6.40956 0.545019 7.59044 1.25221 8.35656L6 13.5" stroke="var(--primary-color)" stroke-linecap="round" style={{transition: 'fill 0.3s'}}/>
                            </svg>
                            Назад
                        </button>

                        <div className='pagination'>
                            {renderPaginationButtons()}
                        </div>

                        <button className='move-button' onClick={() => handlePageChange(currentPage + 1)}>
                            Вперед
                            <svg width="7" height="14" viewBox="0 0 7 14" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M1 0.5L5.74779 5.64344C6.45498 6.40956 6.45498 7.59044 5.74779 8.35656L1 13.5" stroke="var(--primary-color)" stroke-linecap="round" style={{transition: 'fill 0.3s'}}/>
                            </svg>
                        </button>
                    </div>
                )}
            </div>
        </section>
    );
}
