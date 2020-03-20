CREATE TABLE IF NOT EXISTS users (
  id            SERIAL                    NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  email         text                      NOT NULL,
  password      text                      NOT NULL, -- @TODO make this secure
  user_group    integer                   NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS books (
  id            SERIAL                        NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  title         text                          NOT NULL,
  description   text                          ,
  validated_by  integer REFERENCES users (id) NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS book_suggestions (
  id            SERIAL                        NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  suggested_by  integer REFERENCES users (id) NOT NULL,
  approved_by   integer REFERENCES users (id) ,
  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS email_idx ON users (email);
