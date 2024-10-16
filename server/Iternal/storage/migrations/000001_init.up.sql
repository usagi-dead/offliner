CREATE TYPE email_status AS ENUM ('confirmed', 'not confirmed');
-- Создание таблицы Users
CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    hashed_password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    patronymic VARCHAR(100),
    date_of_birth DATE NOT NULL,
    phone_number VARCHAR(15),
    email VARCHAR(100) UNIQUE NOT NULL,
    avatar_url VARCHAR(255) NOT NULL DEFAULT '',
    status_email email_status NOT NULL DEFAULT 'not confirmed',
    gender CHAR(1) CHECK (gender IN ('M', 'F'))
);

-- Создание таблицы Orders
CREATE TABLE Orders (
                        order_id SERIAL PRIMARY KEY,
                        user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
                        order_date DATE NOT NULL,
                        locale VARCHAR(50) NOT NULL,
                        total_price DECIMAL(10, 2) NOT NULL
);

-- Создание таблицы Manufacturer
CREATE TABLE Manufacturer (
                              manufacturer_id SERIAL PRIMARY KEY,
                              name VARCHAR(255) NOT NULL
);

-- Создание таблицы Category
CREATE TABLE Category (
                          category_id SERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL
);

CREATE TABLE Product (
                         product_id SERIAL PRIMARY KEY,
                         manufacturer_id INT REFERENCES Manufacturer(manufacturer_id) ON DELETE SET NULL,
                         category_id INT REFERENCES Category(category_id) ON DELETE SET NULL,
                         name VARCHAR(255) NOT NULL,
                         price DECIMAL(10, 2) NOT NULL,
                         discount DECIMAL(5, 2) CHECK (discount >= 0 AND discount <= 100),
                         discount_price DECIMAL(10, 2)
);

-- Создание триггера для автоматического вычисления discount_price
CREATE OR REPLACE FUNCTION update_discount_price()
RETURNS TRIGGER AS $$
BEGIN
    NEW.discount_price := NEW.price - (NEW.price * NEW.discount / 100);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Применение триггера на вставку и обновление
CREATE TRIGGER calculate_discount_price
    BEFORE INSERT OR UPDATE ON Product
                         FOR EACH ROW EXECUTE FUNCTION update_discount_price();

-- Создание таблицы Order_Items
CREATE TABLE Order_Items (
                             order_id INT REFERENCES Orders(order_id) ON DELETE CASCADE,
                             product_id INT REFERENCES Product(product_id) ON DELETE CASCADE,
                             quantity INT NOT NULL CHECK (quantity > 0),
                             price_at_time DECIMAL(10, 2) NOT NULL,
                             PRIMARY KEY (order_id, product_id)
);

-- Создание таблицы Attributes
CREATE TABLE Attributes (
                            attribute_id SERIAL PRIMARY KEY,
                            category_id INT REFERENCES Category(category_id) ON DELETE CASCADE,
                            attribute_name VARCHAR(255) NOT NULL
);


-- Создание таблицы ProductAttributes
CREATE TABLE ProductAttributes (
                                   product_id INT REFERENCES Product(product_id) ON DELETE CASCADE,
                                   attribute_id INT REFERENCES Attributes(attribute_id) ON DELETE CASCADE,
                                   attribute_value VARCHAR(255) NOT NULL,
                                   PRIMARY KEY (product_id, attribute_id)
);


-- Создание индексов для ускорения поиска
CREATE INDEX idx_users_email ON Users(email);
CREATE INDEX idx_orders_user_id ON Orders(user_id);
CREATE INDEX idx_order_items_product_id ON Order_Items(product_id);
CREATE INDEX idx_product_category_id ON Product(category_id);
CREATE INDEX idx_product_manufacturer_id ON Product(manufacturer_id);
