DROP TABLE IF EXISTS productsInfo;
CREATE TABLE productsInfo (
    seller_id bigint NOT NULL,
    offer_id bigint NOT NULL,
    name varchar(255) NOT NULL,
    price numeric NOT NULL,
    quantity bigint NOT NULL,
    available boolean NOT NULL
);

DROP TABLE IF EXISTS productUploadsTask;
CREATE TABLE productUploadsTask (
    task_id bigserial NOT NULL,
    state varchar(50) NOT NULL
);

DROP TABLE IF EXISTS productTaskStats;
CREATE TABLE productTaskStats (
    task_id bigint NOT NULL,
    products_created bigint,
    products_updated bigint,
    products_deleted bigint,
    rows_with_errors bigint
);