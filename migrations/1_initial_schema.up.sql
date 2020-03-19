CREATE TABLE book_basics (
  id            SERIAL                    NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  title         text                      NOT NULL,
  description   text                      ,
  PRIMARY KEY (id)
);
