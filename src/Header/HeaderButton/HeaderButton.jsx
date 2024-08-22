import React, { useState } from 'react';
import Button from "./Buttons/Button"
import ThemeToggle from "./Buttons/ThemeToggle"
import BusketButton from "./Buttons/BusketButton"

export default function HeaderButton({ link, svg }) {
    const [isDarkTheme, setIsDarkTheme] = useState(false);

    const toggleTheme = () => {
        const newTheme = isDarkTheme ? 'white' : 'black';
        setIsDarkTheme(!isDarkTheme);
        document.documentElement.setAttribute('data-theme', newTheme);
    };

    return (
        <>
            { 
            svg == "theme" ? <ThemeToggle /> :
            svg == "busket" ? <BusketButton link={link} svg={svg} /> : <Button link={link} svg={svg} />
            }  
        </>
    )
}
