import classes from "./ThemeSwitchButton.module.css"
import svgIcons from "../../svgIcons";
import React, { useState } from 'react';

export default function ThemeSwitchButton() {
    const [isDarkTheme, setIsDarkTheme] = useState(false);

    const toggleTheme = () => {
        const newTheme = isDarkTheme ? 'white' : 'black';
        setIsDarkTheme(!isDarkTheme);
        document.documentElement.setAttribute('data-theme', newTheme);
    };

    return (
        <button onClick={toggleTheme} className={classes.themeButton}>
            {svgIcons[isDarkTheme ? 'white' : 'black']}
        </button>
    )
}
