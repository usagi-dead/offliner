CREATE TABLE images (
  image_id SERIAL PRIMARY KEY,
  url TEXT
);

CREATE TABLE productImages (
    product_id INT REFERENCES product(product_id) ON DELETE CASCADE,
    image_id INT REFERENCES images(image_id) ON DELETE CASCADE
);

