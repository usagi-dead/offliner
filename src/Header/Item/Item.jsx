// Item.js
import { Link } from 'react-router-dom';
import classes from "./Item.module.css"

export default function Item({ text, link, svg }) {
    return (
        <Link to={link} className={classes.itemButton}>
            {svg}
            <div className={classes.itemText}>{text}</div>
            {text === 'Корзина' ? <div className={classes.count}>0</div> : null}
        </Link>
    )
}
