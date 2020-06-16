package parking

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	CreateLot(ctx context.Context, maxslotscount uint32) error
	PostPark(ctx context.Context, carreg string, carcolour string) (*Park, error)
	PostUnpark(ctx context.Context, slotnum uint32) error
	GetParks(ctx context.Context) ([]Park, error)
	GetCarRegsByColour(ctx context.Context, carcolour string) ([]Car, error)
	GetSlotsByColour(ctx context.Context, carcolour string) ([]Slot, error)
	GetSlotByCarReg(ctx context.Context, carreg string) (*Slot, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) CreateLot(ctx context.Context, maxslotscount uint32) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO parking_lots(max_slots_count) VALUES($1)`, maxslotscount)
	return err
}

func (r *postgresRepository) PostPark(ctx context.Context, carreg string, carcolour string) (*Park, error) {
	p := &Park{}
	//select recent created parking lot
	row := r.db.QueryRowContext(ctx, `SELECT id, max_slots_count, used_slots_count, next_slot_num FROM parking_lots ORDER BY created_at DESC LIMIT 1)`)
	var ParkingLotID uint32
	var MaxSlotsCount uint32
	var UsedSlotsCount uint32
	var NextSlotNum uint32
	if err := row.Scan(ParkingLotID, MaxSlotsCount, UsedSlotsCount, NextSlotNum); err != nil {
		return nil, err
	}
	if MaxSlotsCount == UsedSlotsCount {
		return nil, ErrParkingFull
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	if _, err := tx.Exec(`INSERT INTO parks(parking_lot_id, slot_num, car_reg, car_colour) VALUES($1, $2, $3, $4)`, ParkingLotID, NextSlotNum, carreg, carcolour); err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `SELECT slot_num FROM parks WHERE parking_lot_id = $1 ORDER BY slot_num ASC)`, ParkingLotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	slots_occupied := make([]uint32, 0, checkCount(rows))
	for rows.Next() {
		var slot uint32
		if err = rows.Scan(&slot); err != nil {
			return nil, err
		}
		slots_occupied = append(slots_occupied, slot)
	}

	var i uint32
	i = 0
	for i < MaxSlotsCount {
		_, found := find(slots_occupied, i)
		if !found {
			NextSlotNum = slots_occupied[i]
			break
		}
	}

	if _, err := tx.Exec(`UPDATE parking_lots SET next_slot_num = $1 WHERE id = $2`, NextSlotNum, ParkingLotID); err != nil {
		return nil, err
	}
	row = r.db.QueryRowContext(ctx, `SELECT slot_num, car_reg, car_colour FROM parks WHERE parking_lot_id = $1)`, ParkingLotID)
	if err := row.Scan(&p.SlotNum, &p.CarReg, &p.CarColour); err != nil {
		return nil, err
	}

	return p, nil
}

func (r *postgresRepository) PostUnpark(ctx context.Context, slotnum uint32) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM parks WHERE slot_num = $1", slotnum)
	return err
}

func (r *postgresRepository) GetParks(ctx context.Context) ([]Park, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id FROM parking_lots ORDER BY created_at DESC LIMIT 1)`)
	var ParkingLotID uint32
	if err := row.Scan(ParkingLotID); err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `SELECT slot_num, car_reg, car_colour FROM parks WHERE parking_lot_id = $1 ORDER BY slot_num ASC)`, ParkingLotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	parks := []Park{}
	for rows.Next() {
		park := &Park{}
		if err = rows.Scan(&park.SlotNum, &park.CarReg, &park.CarColour); err != nil {
			return nil, err
		}
		parks = append(parks, *park)
	}
	return parks, nil
}

func (r *postgresRepository) GetCarRegsByColour(ctx context.Context, carcolour string) ([]Car, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id FROM parking_lots ORDER BY created_at DESC LIMIT 1)`)
	var ParkingLotID uint32
	if err := row.Scan(ParkingLotID); err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `SELECT car_reg FROM parks WHERE parking_lot_id = $1 and car_colour = $2 ORDER BY slot_num ASC)`, ParkingLotID, carcolour)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cars := []Car{}
	for rows.Next() {
		car := &Car{}
		if err = rows.Scan(&car.CarReg); err != nil {
			return nil, err
		}
		cars = append(cars, *car)
	}
	return cars, nil
}

func (r *postgresRepository) GetSlotsByColour(ctx context.Context, carcolour string) ([]Slot, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id FROM parking_lots ORDER BY created_at DESC LIMIT 1)`)
	var ParkingLotID uint32
	if err := row.Scan(ParkingLotID); err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `SELECT slot_num FROM parks WHERE parking_lot_id = $1 and car_colour = $2 ORDER BY slot_num ASC)`, ParkingLotID, carcolour)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	slots := []Slot{}
	for rows.Next() {
		slot := &Slot{}
		if err = rows.Scan(&slot.SlotNum); err != nil {
			return nil, err
		}
		slots = append(slots, *slot)
	}
	return slots, nil
}

func (r *postgresRepository) GetSlotByCarReg(ctx context.Context, carreg string) (*Slot, error) {
	slot := &Slot{}
	row := r.db.QueryRowContext(ctx, `SELECT slot_num FROM parks WHERE car_reg = $1)`, carreg)

	if err := row.Scan(&slot.SlotNum); err != nil {
		return nil, err
	}

	return slot, nil
}
