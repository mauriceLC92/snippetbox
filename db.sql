CREATE DATABASE snippetbox;

-- Take note, timestamptz might be the one I am looking for if I want to store the values by default in UTC
-- CREATE TABLE snippets (
--     id SERIAL,
--     title VARCHAR(100) NOT NULL,
--     content TEXT NOT NULL,
--     created timestamptz NOT NULL,
--     expires timestamptz NOT NULL
-- );
-- https://gist.github.com/alexedwards/dc3145c8e2e6d2fd6cd9
-- If you want a db column to be null
CREATE TABLE snippets (
    id SERIAL,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP NOT NULL,
    expires TIMESTAMP NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);


-- https://x-team.com/blog/automatic-timestamps-with-postgresql/


-- Add some dummy records (which we'll use in the next couple of chapters).
INSERT INTO snippets (title, content, created, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    now(),
    now() + INTERVAL '365 days'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    now(),
    now() + INTERVAL '365 days'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'Fourth autumn morning',
    'Fourth autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    now(),
    now() + INTERVAL '365 days'
);


create user web with encrypted password 'password';

GRANT SELECT, INSERT, UPDATE, DELETE ON snippets TO web;
-- need to do this since sql 8 or something - https://stackoverflow.com/questions/9325017/error-permission-denied-for-sequence-cities-id-seq-using-postgres
GRANT SELECT, INSERT, UPDATE, DELETE ON snippets_id_seq TO web;

SET TIMEZONE='UTC';