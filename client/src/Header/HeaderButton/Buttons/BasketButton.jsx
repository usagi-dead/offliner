import { Link, useLocation } from 'react-router-dom';
import classes from "../HeaderButton.module.css"
import svgIcons from "../../../svgIcons";

export default function BasketButton({ link, svg, svgFilled }) {
    const location = useLocation();

    return (
        <Link to={link} className={classes.itemButton}>
            {
                location.pathname.startsWith(`${link}`) ?
                svgIcons[svgFilled] :
                svgIcons[svg]
            }
            <div className={classes.count}>0</div>
        </Link>
    )
}
