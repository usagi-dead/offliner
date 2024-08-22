import { Link } from 'react-router-dom';
import classes from "../HeaderButton.module.css"
import svgIcons from "../../../svgIcons";

export default function BusketButton({ link, svg }) {
    return (
        <Link to={link} className={classes.itemButton}>
            {svgIcons[svg]}
            <div className={classes.count}>0</div>
        </Link>
    )
}
