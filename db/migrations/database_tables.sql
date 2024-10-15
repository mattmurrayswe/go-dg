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
    site VARCHAR(255),
    logo_url VARCHAR(255),
    UNIQUE (brand_name)
);

CREATE TABLE dg_brands (
    id SERIAL PRIMARY KEY,
    brand_name VARCHAR(255),
    site VARCHAR(255),
    logo_url VARCHAR(255),
    UNIQUE (brand_name)
);

CREATE TABLE models (
    id SERIAL PRIMARY KEY,
    brand_name VARCHAR(255),
    model_name VARCHAR(255),
    UNIQUE (brand_name, model_name)
);

CREATE TABLE versions (
    id SERIAL PRIMARY KEY,
    brand_name VARCHAR(255),
    model_name VARCHAR(255),
    version_name VARCHAR(255),
    UNIQUE(brand_name, model_name, version_name)
);

CREATE TABLE wishlist (
    id SERIAL PRIMARY KEY,
    brand VARCHAR(255),
    model VARCHAR(255),
    year VARCHAR(255)
)

-- SELECT setval('models_id_seq', 1, false);
