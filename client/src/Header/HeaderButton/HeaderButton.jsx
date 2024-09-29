import React from 'react';
import Button from "./Buttons/Button"
import ThemeToggle from "./Buttons/ThemeToggle"
import BusketButton from "./Buttons/BusketButton"

export default function HeaderButton({ link, svg }) {

    return (
        <>
            { 
            svg == "theme" ? <ThemeToggle /> :
            svg == "busket" ? <BusketButton link={link} svg={svg} /> : <Button link={link} svg={svg} />
            }  
        </>
    )
}
