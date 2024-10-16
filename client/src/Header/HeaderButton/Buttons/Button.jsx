import { Link, useLocation } from 'react-router-dom';
import classes from "../HeaderButton.module.css"
import svgIcons from "../../../svgIcons";

export default function Button({ link, svg, svgFilled }) {
    const location = useLocation();

    return (
        <Link to={link} className={classes.itemButton}>
            {
                location.pathname.startsWith(`${link}`) ?
                svgIcons[svgFilled] :
                svgIcons[svg]
            }
        </Link>
    )
}
