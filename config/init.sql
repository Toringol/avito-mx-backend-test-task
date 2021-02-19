DROP TABLE IF EXISTS productsInfo;
CREATE TABLE productsInfo (
    seller_id bigint NOT NULL,
    offer_id bigint NOT NULL,
    name varchar(255) NOT NULL,
    price numeric NOT NULL,
    quantity bigint NOT NULL,
    available boolean NOT NULL
);