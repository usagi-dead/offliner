import { Link } from 'react-router-dom';
import classes from "./HeaderButton.module.css"
import svgIcons from "../../svgIcons";

export default function Item({ link, svg, isBusket }) {
    return (
        <Link to={link} className={isBusket == 1 ? classes.busketButton : classes.itemButton}>
            {svgIcons[svg]}
            {isBusket == 1 ? <div className={classes.count}>0</div> : <></>}
        </Link>
    )
}
