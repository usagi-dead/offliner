import React from 'react';
import Button from "./Buttons/Button"
import ThemeToggle from "./Buttons/ThemeToggle"
import BasketButton from "./Buttons/BasketButton"

export default function HeaderButton({ link, svg, svgFilled }) {
    return (
        <>
            { 
            svg == "theme" ? <ThemeToggle /> :
            svg == "basket" ? <BasketButton link={link} svg={svg} svgFilled={svgFilled} /> : <Button link={link} svg={svg} svgFilled={svgFilled} />
            }  
        </>
    )
}
