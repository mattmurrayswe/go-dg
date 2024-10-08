CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE tech_options (
    id SERIAL PRIMARY KEY,
    tech_name VARCHAR(255) NOT NULL UNIQUE,
    image_url VARCHAR(255) NOT NULL
);

CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    owner_id INT,
    brand VARCHAR(255),
    model VARCHAR(255),
    year INT,
    card_price INT,
    project_name VARCHAR(255),
    photo VARCHAR(255),
    horse_powers INT,
    dgp INT,
    rarity VARCHAR(255)
);

CREATE TABLE brands (
    id SERIAL PRIMARY KEY,
    brand_name VARCHAR(255),
);

CREATE TABLE dg_brands (
    id SERIAL PRIMARY KEY,
    brand_name VARCHAR(255),
);

CREATE TABLE models (
    id SERIAL PRIMARY KEY,
    brand_name VARCHAR(255),
    model_name VARCHAR(255)
);

-- CREATE TABLE models (
--     id SERIAL PRIMARY KEY,
--     brand_name VARCHAR(255),
--     model_name VARCHAR(255),
--     UNIQUE(brand_name, model_name)
-- );

ALTER table dg_brands 
ADD COLUMN site VARCHAR(255);

ALTER table brands 
ADD COLUMN site VARCHAR(255);

ALTER table brands 
ADD COLUMN logo_url VARCHAR(255);

ALTER TABLE models
ADD CONSTRAINT unique_brand_model UNIQUE (brand_name, model_name);

-- SELECT setval('models_id_seq', 1, false);

SELECT models.model_name, models.brand_name, dg_brands.site
FROM models
INNER JOIN dg_brands ON models.brand_name = dg_brands.brand_name;
