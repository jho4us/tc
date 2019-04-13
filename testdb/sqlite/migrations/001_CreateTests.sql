-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE tests (
  id                       char(38) NOT NULL,
  name                     varchar(128) NOT NULL,
  created_at               timestamp DEFAULT CURRENT_TIMESTAMP,
  modified_at              timestamp DEFAULT CURRENT_TIMESTAMP,
  content                  TEXT,
  PRIMARY KEY(id)
);

CREATE INDEX ni ON tests(name);

CREATE TRIGGER UpdateLastTime AFTER UPDATE OF name, content ON tests FOR EACH ROW WHEN NEW.modified_at <= OLD.modified_at BEGIN update tests set modified_at=CURRENT_TIMESTAMP where id=OLD.id; END;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE tests;

