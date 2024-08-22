import { useRef } from 'react';
import classes from './Card.module.css';
import { Link } from 'react-router-dom';

export default function Card({ title, cardImage, urlProduct }) {
    const cardRef = useRef(null);

    const handleMouseMove = (e) => {
        const card = cardRef.current;
        const rect = card.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const y = e.clientY - rect.top;
        const centerX = rect.width / 2;
        const centerY = rect.height / 2;
        const rotateX = ((y - centerY) / centerY) * 10;
        const rotateY = ((x - centerX) / centerX) * -10; 
        card.style.transform = `perspective(1000px) rotateX(${rotateX}deg) rotateY(${rotateY}deg)`;
    };

    const handleMouseLeave = () => {
        const card = cardRef.current;
        card.style.transform = 'perspective(1000px) rotateX(0deg) rotateY(0deg)';
    };

    return (
        <Link to={`/${urlProduct}/1`} className={classes.cardLink}>
            <button
                className={classes.card}
                onMouseMove={handleMouseMove}
                onMouseLeave={handleMouseLeave}
                ref={cardRef}
            >
                <img src={cardImage} alt={cardImage} />
                <h2>{title}</h2>
            </button>
        </Link>
    );
}
