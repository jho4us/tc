-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE tests (
  id                       char(38) NOT NULL,
  name                     varchar(128) NOT NULL,
  created_at               timestamp DEFAULT now(),
  modified_at              timestamp DEFAULT now() ON UPDATE now(),
  content                  LONGTEXT,
  PRIMARY KEY(id),
  KEY (`name`),
  FULLTEXT KEY `ft1` (`name`,`content`),
  FULLTEXT KEY `ft2` (`content`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE tests;
