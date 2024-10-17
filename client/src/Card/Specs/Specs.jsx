import React from 'react';
import classes from './Specs.module.css';

export default function Specs({ specs, textSize, lineNum = 3 }) {
    const DISPLAYED_SPECS = [
        "Технические характеристики: Видеопамять",
        "Технические характеристики: Тип видеопамяти",
        "Технические характеристики: Разъёмы питания",
        "Технические характеристики: Энергопотребление",
        "Технические характеристики: Толщина системы охлаждения",
        "Технические характеристики: Функциональные особенности",
    ];

    const displayedSpecs = Object.keys(specs).filter(key => DISPLAYED_SPECS.includes(key))
        .map(key => {
            const value = specs[key];
            return { key: key.split(': ')[1], value };
        });

    return (
        <ul className={classes.specList} style={{ height: (textSize * (lineNum + 1) + 8) + "px" }}>
            {displayedSpecs.map((spec, index) => (
                <li key={index} className={classes.specItem} style={{ fontSize: textSize + "px" }}>
                    {spec.value}
                </li>
            ))}
        </ul>
    );
}
