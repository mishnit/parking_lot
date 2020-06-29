package parking

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	CreateLot(ctx context.Context, maxslotscount uint32) error
	PostPark(ctx context.Context, carreg string, carcolour string) (*Park, error)
	PostUnpark(ctx context.Context, slotnum uint32) error
	GetParks(ctx context.Context) ([]Park, error)
	GetCarRegsByColour(ctx context.Context, carcolour string) ([]string, error)
	GetSlotsByColour(ctx context.Context, carcolour string) ([]uint32, error)
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

	_, err := r.db.ExecContext(ctx, "DELETE FROM parks")

	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, "DELETE FROM parking_lots")

	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `INSERT INTO parking_lots(max_slots_count) VALUES($1)`, maxslotscount)
	return err
}

func (r *postgresRepository) PostPark(ctx context.Context, carreg string, carcolour string) (*Park, error) {
	p := &Park{}

	if !RegexCarNumber.MatchString(carreg) {
		return nil, ErrInvalidCarNumber
	}
	//select recent created parking lot
	row := r.db.QueryRowContext(ctx, `SELECT id, max_slots_count, used_slots_count, next_slot_num FROM parking_lots ORDER BY created_at DESC LIMIT 1`)
	var ParkingLotID uint32
	var MaxSlotsCount uint32
	var UsedSlotsCount uint32
	var CurrSlotNum uint32
	if err := row.Scan(&ParkingLotID, &MaxSlotsCount, &UsedSlotsCount, &CurrSlotNum); err != nil {
		//log.Println("repo debug: PostPark-> scan parking_lot")
		return nil, ErrNoLotFound
	}

	// is parking lot full?
	if MaxSlotsCount == UsedSlotsCount {
		return nil, ErrParkingFull
	}

	//log.Println("repo debug: PostPark-> before nextslotnum: ", CurrSlotNum)

	// every insert will update used_slots_count and used_slots
	// now we need to update next available slot

	rows, err := r.db.QueryContext(ctx, `SELECT slot_num FROM parks WHERE parking_lot_id = $1 ORDER BY slot_num ASC`, ParkingLotID)
	if err != nil {
		//log.Println("repo debug: GetParks-> select slot_num")
		return nil, err
	}
	defer rows.Close()

	UsedSlots := []uint32{}
	for rows.Next() {
		var UsedSlot uint32
		if err = rows.Scan(&UsedSlot); err != nil {
			//log.Println("repo debug: GetParks-> scan park")
			return nil, err
		}
		UsedSlots = append(UsedSlots, UsedSlot)
	}

	UsedSlots = append(UsedSlots, CurrSlotNum) //precompute

	NextSlotNum := nextslot(MaxSlotsCount, UsedSlots)

	//log.Println("repo debug: PostPark-> after nextslotnum: ", NextSlotNum)

	/// transaction begin

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tx.ExecContext(ctx, `INSERT INTO parks(parking_lot_id, slot_num, car_reg, car_colour) VALUES($1, $2, $3, $4)`, ParkingLotID, CurrSlotNum, carreg, carcolour); err != nil {
		//log.Println("repo debug: PostPark-> insert park")
		tx.Rollback()
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `UPDATE parking_lots SET next_slot_num = $1 WHERE id = $2`, NextSlotNum, ParkingLotID); err != nil {
		//log.Println("repo debug: PostPark-> update next_slot_num")
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	/// transaction end

	row = r.db.QueryRowContext(ctx, `SELECT slot_num, car_reg, car_colour FROM parks WHERE parking_lot_id = $1 and slot_num = $2`, ParkingLotID, CurrSlotNum)
	if err := row.Scan(&p.SlotNum, &p.CarReg, &p.CarColour); err != nil {
		//log.Println("repo debug: PostPark-> select slot_num, car_reg, car_colour")
		return nil, err
	}

	return p, nil
}

func (r *postgresRepository) PostUnpark(ctx context.Context, slotnum uint32) error {
	var ParkingLotID uint32
	var UsedSlotsCount uint32
	var MaxSlotsCount uint32
	var CurrSlotNum uint32

	row := r.db.QueryRowContext(ctx, `SELECT id, max_slots_count, used_slots_count, next_slot_num FROM parking_lots ORDER BY created_at DESC LIMIT 1`)
	if err := row.Scan(&ParkingLotID, &MaxSlotsCount, &UsedSlotsCount, &CurrSlotNum); err != nil {
		//log.Println("repo debug: PostUnpark-> scan parking_lot")
		return ErrNoLotFound
	}

	if slotnum > MaxSlotsCount {
		return ErrInvalidSlot
	}

	// is parking slot already empty?
	var slot_num_exists_for_lot bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM parks WHERE slot_num = $1 and parking_lot_id = $2)`, slotnum, ParkingLotID).Scan(&slot_num_exists_for_lot)
	if err != nil {
		return fmt.Errorf("could not query select parks existence: %w", err)
	}

	if !slot_num_exists_for_lot {
		return ErrParking
	}

	_, err = r.db.ExecContext(ctx, `DELETE FROM parks WHERE slot_num = $1 and parking_lot_id = $2`, slotnum, ParkingLotID)
	if err != nil {
		return err
	}
	//log.Println("repo debug: PostUnpark-> before nextslotnum: ", CurrSlotNum)
	// every delete will update used_slots_count and used_slots
	// now we need to update next available slot
	rows, err := r.db.QueryContext(ctx, `SELECT slot_num FROM parks WHERE parking_lot_id = $1 ORDER BY slot_num ASC`, ParkingLotID)
	if err != nil {
		//log.Println("repo debug: GetParks-> select slot_num")
		return err
	}
	defer rows.Close()

	UsedSlots := []uint32{}
	for rows.Next() {
		var UsedSlot uint32
		if err = rows.Scan(&UsedSlot); err != nil {
			//log.Println("repo debug: GetParks-> scan park")
			return err
		}
		UsedSlots = append(UsedSlots, UsedSlot)
	}
	NextSlotNum := nextslot(MaxSlotsCount, UsedSlots)
	//log.Println("repo debug: PostUnpark-> after nextslotnum: ", NextSlotNum)

	if _, err := r.db.ExecContext(ctx, `UPDATE parking_lots SET next_slot_num = $1 WHERE id = $2`, NextSlotNum, ParkingLotID); err != nil {
		//log.Println("repo debug: PostUnpark-> update next_slot_num")
		return err
	}

	return nil
}

func (r *postgresRepository) GetParks(ctx context.Context) ([]Park, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id FROM parking_lots ORDER BY created_at DESC LIMIT 1`)
	var ParkingLotID uint32
	if err := row.Scan(&ParkingLotID); err != nil {
		//log.Println("repo debug: GetParks-> select id")
		return nil, ErrNoLotFound
	}

	rows, err := r.db.QueryContext(ctx, `SELECT slot_num, car_reg, car_colour FROM parks WHERE parking_lot_id = $1 ORDER BY slot_num ASC`, ParkingLotID)
	if err != nil {
		//log.Println("repo debug: GetParks-> select slot_num car_reg car_colour")
		return nil, err
	}
	defer rows.Close()

	parks := []Park{}
	for rows.Next() {
		park := &Park{}
		if err = rows.Scan(&park.SlotNum, &park.CarReg, &park.CarColour); err != nil {
			//log.Println("repo debug: GetParks-> scan park")
			return nil, ErrNotFound
		}
		parks = append(parks, *park)
	}
	return parks, nil
}

func (r *postgresRepository) GetCarRegsByColour(ctx context.Context, carcolour string) ([]string, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id FROM parking_lots ORDER BY created_at DESC LIMIT 1`)
	var ParkingLotID uint32
	if err := row.Scan(&ParkingLotID); err != nil {
		//log.Println("repo debug: GetCarRegsByColour-> select id")
		return nil, ErrNoLotFound
	}

	rows, err := r.db.QueryContext(ctx, `SELECT car_reg FROM parks WHERE parking_lot_id = $1 and car_colour = $2 ORDER BY slot_num ASC`, ParkingLotID, carcolour)
	if err != nil {
		//log.Println("repo debug: GetCarRegsByColour-> select car_reg")
		return nil, err
	}
	defer rows.Close()

	cars := []string{}
	for rows.Next() {
		var CarReg string
		if err = rows.Scan(&CarReg); err != nil {
			//log.Println("repo debug: GetCarRegsByColour-> scan car")
			return nil, ErrNotFound
		}
		cars = append(cars, CarReg)
	}
	return cars, nil
}

func (r *postgresRepository) GetSlotsByColour(ctx context.Context, carcolour string) ([]uint32, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id FROM parking_lots ORDER BY created_at DESC LIMIT 1`)
	var ParkingLotID uint32
	if err := row.Scan(&ParkingLotID); err != nil {
		//log.Println("repo debug: GetSlotsByColour-> select id")
		return nil, ErrNoLotFound
	}

	rows, err := r.db.QueryContext(ctx, `SELECT slot_num FROM parks WHERE parking_lot_id = $1 and car_colour = $2 ORDER BY slot_num ASC`, ParkingLotID, carcolour)
	if err != nil {
		//log.Println("repo debug: GetSlotsByColour-> select slot_num")
		return nil, err
	}
	defer rows.Close()

	slots := []uint32{}
	for rows.Next() {
		var SlotNum uint32
		if err = rows.Scan(&SlotNum); err != nil {
			//log.Println("repo debug: GetSlotsByColour-> scan slot")
			return nil, ErrNotFound
		}
		slots = append(slots, SlotNum)
	}
	return slots, nil
}

func (r *postgresRepository) GetSlotByCarReg(ctx context.Context, carreg string) (*Slot, error) {
	if !RegexCarNumber.MatchString(carreg) {
		return nil, ErrInvalidCarNumber
	}
	slot := &Slot{}
	row := r.db.QueryRowContext(ctx, `SELECT slot_num FROM parks WHERE car_reg = $1`, carreg)

	if err := row.Scan(&slot.SlotNum); err != nil {
		//log.Println("repo debug: GetSlotByCarReg")
		return nil, ErrNotFound
	}

	return slot, nil
}
