import React from 'react';
import "./Pagination.css";
import svgIcons from "../../svgIcons";

export default function Pagination({ currentPage, totalPages, onPageChange }) {
    const maxPagesToShow = 9;

    const handlePageChange = (page) => {
        if (page > 0 && page <= totalPages) {
            onPageChange(page);
        }
    };

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
        <div className="pagination-container">
            <button className='move-button' onClick={() => handlePageChange(currentPage - 1)}>
                <span className='left-arrow'>{svgIcons["smallArrow"]}</span>
                Назад
            </button>

            <div className='pagination'>
                {renderPaginationButtons()}
            </div>

            <button className='move-button' onClick={() => handlePageChange(currentPage + 1)}>
                Вперед
                {svgIcons["smallArrow"]}
            </button>
        </div>
    );
}
