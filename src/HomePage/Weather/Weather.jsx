import React, { useState, useEffect } from 'react';
import "./Weather.css";

const API_KEY_WEATHER = 'd8da04f2fbeae7b5378ae6e65264ae5b';

const weatherImages = {
    '01d': 'clear_day.svg',
    '01n': 'clear_night.svg',
    '02d': 'few_clouds_day.svg',
    '02n': 'few_clouds_night.svg',
    '03d': 'clouds.svg',
    '03n': 'clouds.svg',
    '04d': 'clouds.svg',
    '04n': 'clouds.svg',
    '09d': 'rain.svg',
    '09n': 'rain.svg',
    '10d': 'rain.svg',
    '10n': 'rain.svg',
    '11d': 'thunderstorm.svg',
    '11n': 'thunderstorm.svg',
    '13d': 'snow.svg',
    '13n': 'snow.svg',
    '50d': 'mist.svg',
    '50n': 'mist.svg',
};

const weatherPhrases = {
    '01d': 'Сегодня солнечно! Отличный день, чтобы сходить и забрать свои комплектующие.',
    '01n': 'Ясная ночь. Прекрасное время забрать свои комплектующие на такси.',
    '02d': 'Немного облаков, но солнце светит. Пора забрать свои комплектующие!',
    '02n': 'Легкая облачность. Заберите свои комплектующие сейчас, пока тихо.',
    '03d': 'Облачно, но без дождя. Идеально для прогулки до магазина за комплектующими.',
    '03n': 'Облака заволокли небо. Заезжайте за комплектующими по пути домой.',
    '04d': 'Пасмурно, но не унывайте. Самое время купить новые комплектующие.',
    '04n': 'Пасмурная ночь. Заберите свои комплектующие на такси.',
    '09d': 'Дождь идет. Возьмите зонт и заберите свои комплектующие.',
    '09n': 'Ночная непогода за окном. Закажите комплектующие с доставкой.',
    '10d': 'Дождливый день. Не забудьте зонт и заберите свои комплектующие.',
    '10n': 'Дождь ночью. Закажите такси и заберите свои комплектующие.',
    '11d': 'Гроза на горизонте! Закажите комплектующие онлайн.',
    '11n': 'Гроза ночью. Закажите комплектующие с доставкой на дом.',
    '13d': 'Снег падает! Пора заехать за комплектующими и сделать свои дела.',
    '13n': 'Снежная ночь. Закажите такси и заберите свои комплектующие.',
    '50d': 'Туман окутал всё. Будьте осторожны на дорогах, забирая комплектующие.',
    '50n': 'Туманная ночь. Видимость ограничена, но комплектующие нужно забрать!',
};

export default function Weather() {
    const [weatherData, setWeatherData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchWeatherData = async () => {
            try {
                const response = await fetch(`https://api.openweathermap.org/data/2.5/weather?q=Москва&units=metric&lang=ru&appid=${API_KEY_WEATHER}`);
                const data = await response.json();
                if (response.ok) {
                    setWeatherData(data);
                } else {
                    throw new Error(data.message);
                }
            } catch (error) {
                setError(error.message);
            } finally {
                setLoading(false);
            }
        };

        fetchWeatherData();
    }, []);

    if (loading) {
        return (
            <section className="weather-container">
                <h1 className="title">Погода сейчас</h1>
                <div className="weather-data">
                    <h1 className="date">Загрузка . . .</h1>
                </div>
            </section>
        );
    }

    if (error) {
        return <div>Ошибка: {error}</div>;
    }

    const formatDescription = (description) => {
        return description.charAt(0).toUpperCase() + description.slice(1);
    };

    const formatDate = () => {
        const options = { weekday: 'long', day: 'numeric', month: 'long' };
        return new Date().toLocaleDateString('ru-RU', options);
    };

    const getWeatherImage = (icon) => {
        const imgName = weatherImages[icon];
        return "/Weather/" + imgName;
    };

    const getWeatherPhrase = (icon) => {
        return weatherPhrases[icon] || "Такую погоду мы еще не встречали!";
    };

    return (
        <section className="weather-container">
            <h1 className="title">Погода сейчас</h1>
            <div className="weather-data">
                <h3 className="date">{formatDate()}</h3>
                <h1 className="degrees">{Math.round(weatherData.main.temp)}°</h1>
                <h3 className="weather-type">{formatDescription(weatherData.weather[0].description)}</h3>
            </div>
            
            <div className="quote">
                <div className="circle">
                    <h2 className="hooks">“</h2>
                </div>
                <div className="line" />
                <p className="quote-text">{getWeatherPhrase(weatherData.weather[0].icon)}</p>
            </div>
            <img src={getWeatherImage(weatherData.weather[0].icon)} alt={formatDescription(weatherData.weather[0].description)} className="weather-img" />
        </section>
    );
}
