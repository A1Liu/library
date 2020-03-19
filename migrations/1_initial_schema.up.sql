CREATE TABLE book_basics (
  id            SERIAL                    NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  title         text                      NOT NULL,
  description   text                      ,
  PRIMARY KEY (id)
);

create table users (
  id            SERIAL                    NOT NULL,
  email         text                      NOT NULL,
  password      text                      NOT NULL, -- @TODO make this secure
  permissions   integer                   NOT NULL,
  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX email_idx ON users (email);
