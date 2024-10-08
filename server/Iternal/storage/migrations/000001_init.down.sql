DROP  TABLE Order_Items;
DROP  TABLE ProductAttributes;
DROP  TABLE Attributes;
DROP  TABLE Product;
DROP  TABLE Manufacturer;
DROP  TABLE Category;
DROP  TABLE Orders;
DROP  TABLE Users;

DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_order_items_product_id;
DROP INDEX IF EXISTS idx_product_category_id;
DROP INDEX IF EXISTS idx_product_manufacturer_id;