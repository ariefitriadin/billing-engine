-- migrate:up
CREATE TABLE borrowers (
    LIKE template_table INCLUDING ALL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) NOT NULL
);

-- migrate:down
DROP TABLE borrowers;