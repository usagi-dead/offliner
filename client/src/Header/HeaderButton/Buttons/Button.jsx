import { Link } from 'react-router-dom';
import classes from "../HeaderButton.module.css"
import svgIcons from "../../../svgIcons";

export default function Button({ link, svg }) {
    return (
        <Link to={link} className={classes.itemButton}>
            {svgIcons[svg]}
        </Link>
    )
}
