-- migrate:up
CREATE TABLE template_table (
    id SERIAL PRIMARY KEY,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deletedAt TIMESTAMP
);

-- migrate:down
DROP TABLE template_table;