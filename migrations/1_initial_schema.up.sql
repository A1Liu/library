CREATE TABLE IF NOT EXISTS users (
  id            SERIAL                                NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE              NOT NULL DEFAULT CURRENT_TIMESTAMP,
  username      varchar(16)                           NOT NULL UNIQUE,
  email         text                                  NOT NULL UNIQUE,
  password      text                                  NOT NULL, -- @TODO make this secure
  user_group    integer                               NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS tokens (
  id            SERIAL                                NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE              NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expires_at    TIMESTAMP WITH TIME ZONE              NOT NULL,
  user_id       integer REFERENCES users (id)         NOT NULL,
  value         char(128)                             NOT NULL UNIQUE,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS images (
  id            SERIAL                                NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE              NOT NULL DEFAULT CURRENT_TIMESTAMP,
  extension     varchar(8)                            NOT NULL,
  data          BYTEA                                 NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS permissions (
  id            SERIAL                                NOT NULL,
  given_at      TIMESTAMP WITH TIME ZONE              NOT NULL DEFAULT CURRENT_TIMESTAMP,
  given_to      integer REFERENCES users (id)         NOT NULL,
  authorized_by integer REFERENCES users (id)         NOT NULL,
  permission_to integer                               NOT NULL,
  metadata      integer                               NOT NULL, -- Could be ref to anything
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS authors (
  id            SERIAL                                NOT NULL,
  suggested_at  TIMESTAMP WITH TIME ZONE              NOT NULL DEFAULT CURRENT_TIMESTAMP,
  suggested_by  integer REFERENCES users (id)         , -- null means suggested anonymously
  validated_at  TIMESTAMP WITH TIME ZONE              DEFAULT CURRENT_TIMESTAMP,
  validated_by  integer REFERENCES users (id)         ,
  first_name    varchar(32)                           NOT NULL,
  last_name     varchar(32)                           NOT NULL,
  image         integer REFERENCES images (id)        ,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS books (
  id            SERIAL                                NOT NULL,
  suggested_at  TIMESTAMP WITH TIME ZONE              NOT NULL DEFAULT CURRENT_TIMESTAMP,
  suggested_by  integer REFERENCES users (id)         , -- null means suggested anonymously
  validated_at  TIMESTAMP WITH TIME ZONE              DEFAULT CURRENT_TIMESTAMP,
  validated_by  integer REFERENCES users (id)         ,
  title         text                                  NOT NULL,
  description   text                                  NOT NULL,
  cover         integer REFERENCES images (id)        ,
  PRIMARY KEY (id)
);

-- missing written_by means written anonymously
CREATE TABLE IF NOT EXISTS written_by (
  id            SERIAL                                NOT NULL,
  suggested_at  TIMESTAMP WITH TIME ZONE              NOT NULL DEFAULT CURRENT_TIMESTAMP,
  suggested_by  integer REFERENCES users (id)         ,
  validated_at  TIMESTAMP WITH TIME ZONE              DEFAULT CURRENT_TIMESTAMP,
  validated_by  integer REFERENCES users (id)         ,
  author_id     integer REFERENCES authors (id)       NOT NULL,
  book_id       integer REFERENCES books (id)         NOT NULL,
  PRIMARY KEY (id)
);
