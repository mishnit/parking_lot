CREATE DATABASE IF NOT EXISTS parking;

CREATE TABLE IF NOT EXISTS park_lots (
  id SERIAL PRIMARY KEY,
  max_slot_num INT NOT NULL,
  used_slot_count INT NOT NULL DEFAULT 0 CHECK (used_slot_count >= 0), CHECK (used_slot_count <= max_slot_num)
);

CREATE TABLE IF NOT EXISTS parks (
  park_lot_id INT REFERENCES park_lots, --insertion allowed at latest lot only i.e. max(park_lot_id)
  slot_num INT NOT NULL UNIQUE,
  car_num VARCHAR NOT NULL UNIQUE,
  car_color VARCHAR NOT NULL,
  PRIMARY KEY (park_lot_id, slot_num)
);

CREATE FUNCTION inc_used_slot() RETURNS TRIGGER
AS $_$
BEGIN
UPDATE park_lots SET used_slot_count = used_slot_count + 1 WHERE id = NEW.park_lot_id;
RETURN NEW;
END $_$
LANGUAGE 'plpgsql';

CREATE FUNCTION dec_used_slot() RETURNS TRIGGER
AS $_$
BEGIN
UPDATE park_lots SET used_slot_count = used_slot_count - 1  WHERE id = NEW.park_lot_id;
RETURN NEW;
END $_$
LANGUAGE 'plpgsql';

DROP TRIGGER IF EXISTS inc_used_slot_count ON parks;

CREATE TRIGGER inc_used_slot_count -- trigger to increment used_slot_count in park_lots for inserts in parks
  AFTER INSERT ON parks
  FOR EACH ROW
  EXECUTE PROCEDURE inc_used_slot();

DROP TRIGGER IF EXISTS dec_used_slot_count ON parks;

CREATE TRIGGER dec_used_slot_count -- trigger to decrement used_slot_count in park_lots for delete in parks
  AFTER DELETE ON parks
  FOR EACH ROW
  EXECUTE PROCEDURE dec_used_slot();
