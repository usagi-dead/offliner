import classes from "../HeaderButton.module.css"
import svgIcons from "../../../svgIcons";
import React, { useState } from 'react';

export default function ThemeToggle() {
    const [isDarkTheme, setIsDarkTheme] = useState(false);

    const toggleTheme = () => {
        const newTheme = isDarkTheme ? 'white' : 'black';
        setIsDarkTheme(!isDarkTheme);
        document.documentElement.setAttribute('data-theme', newTheme);
    };

    return (
        <button onClick={toggleTheme} className={classes.itemButton}>
            {svgIcons[isDarkTheme ? 'white' : 'black']}
        </button>
    )
}
