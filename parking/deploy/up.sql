CREATE DATABASE IF NOT EXISTS parking;

CREATE TABLE IF NOT EXISTS parking_lots (
  id SERIAL PRIMARY KEY,
  max_slots_count INT NOT NULL,
  used_slots_count INT NOT NULL DEFAULT 0 CHECK (used_slots_count >= 0), CHECK (used_slots_count <= max_slots_count),
  next_slot_num INT NOT NULL DEFAULT 1 CHECK (next_slot_num >= 1), CHECK (next_slot_num <= max_slots_count+1),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS parks (
  parking_lot_id INT REFERENCES parking_lots,
  slot_num INT NOT NULL UNIQUE,
  car_reg VARCHAR NOT NULL UNIQUE,
  car_colour VARCHAR NOT NULL,
  PRIMARY KEY (park_lot_id, slot_num)
);

CREATE FUNCTION inc_used_slot() RETURNS TRIGGER
AS $_$
BEGIN
UPDATE parking_lots SET used_slots_count = used_slots_count + 1 WHERE id = NEW.park_lot_id;
RETURN NEW;
END $_$
LANGUAGE 'plpgsql';

CREATE FUNCTION dec_used_slot() RETURNS TRIGGER
AS $_$
BEGIN
UPDATE parking_lots SET used_slots_count = used_slots_count - 1  WHERE id = OLD.park_lot_id;
RETURN NEW;
END $_$
LANGUAGE 'plpgsql';



DROP TRIGGER IF EXISTS inc_used_slots_count ON parks;

CREATE TRIGGER inc_used_slots_count -- trigger to increment used_slots_count in parking_lots for inserts in parks
  AFTER INSERT ON parks
  FOR EACH ROW
  EXECUTE PROCEDURE inc_used_slot();

DROP TRIGGER IF EXISTS dec_used_slots_count ON parks;

CREATE TRIGGER dec_used_slots_count -- trigger to decrement used_slots_count in parking_lots for delete in parks
  AFTER DELETE ON parks
  FOR EACH ROW
  EXECUTE PROCEDURE dec_used_slot();
