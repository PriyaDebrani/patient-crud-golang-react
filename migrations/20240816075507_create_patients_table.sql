-- +goose Up
CREATE TABLE patients (
    id int NOT NULL,
    name varchar(255),
    disease varchar(255),
    address  varchar(255),
    date int,
    month int,
    year int,
    phone int,
    PRIMARY KEY(id)
);

-- +goose Down
DROP table patients;
