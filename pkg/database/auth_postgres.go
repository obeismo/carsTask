package database

import (
	"database/sql"
	"em_test/pkg/model"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func AddNewCar(newCar model.Car) (int, error) {
	db, err := InitDB()
	if err != nil {
		logrus.Errorf("failed to initialize to db: %s", err.Error())
		return 0, nil
	}

	var ownerID int
	err = db.QueryRow("SELECT id FROM people WHERE id = $1", newCar.Owner.ID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			row := db.QueryRow("INSERT INTO people (name, surname, patronymic) VALUES ($1, $2, $3) RETURNING id",
				newCar.Owner.Name, newCar.Owner.Surname, newCar.Owner.Patronymic)
			err = row.Scan(&ownerID)
			if err != nil {
				logrus.Errorf("failed to insert data in people table: %s", err.Error())
				return 0, err
			}
		}
	}

	_, err = db.Exec("INSERT INTO car (reg_num, mark, model, year, owner) VALUES ($1, $2, $3, $4, $5)",
		newCar.RegNum, newCar.Mark, newCar.Model, newCar.Year, ownerID)
	if err != nil {
		logrus.Errorf("failed to insert data in car table: %s", err.Error())
		return 0, err
	}

	var carID int
	err = db.QueryRow("SELECT id FROM car WHERE reg_num = $1", newCar.RegNum).Scan(&carID)
	if err != nil {
		return 0, err
	}

	return carID, nil
}

func GetCarsByFilter(filter map[string]string, offset int, limit int) ([]model.Car, error) {
	db, err := InitDB()
	if err != nil {
		logrus.Errorf("failed to initialize to db: %s", err.Error())
		return []model.Car{}, err
	}

	query := "SELECT * FROM car WHERE 1=1"

	if filter["reg_num"] != "" {
		query += " AND reg_num = '" + filter["reg_num"] + "'"
	}
	if filter["mark"] != "" {
		query += " AND mark = '" + filter["mark"] + "'"
	}
	if filter["model"] != "" {
		query += " AND model = '" + filter["model"] + "'"
	}
	if filter["year"] != "" {
		query += " AND year = '" + filter["year"] + "'"
	}

	offset = (offset - 1) * limit
	limitStr, offsetStr := strconv.Itoa(limit), strconv.Itoa(offset)
	query += " LIMIT " + limitStr + "OFFSET " + offsetStr

	rows, err := db.Query(query)
	if err != nil {
		return []model.Car{}, err
	}
	defer rows.Close()

	var cars []model.Car
	var car model.Car
	for rows.Next() {
		err := rows.Scan(&car.ID, &car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.ID)
		if err != nil {
			return []model.Car{}, err
		}
		if car.Owner.ID > 0 {
			rowsOwner, err := db.Query("SELECT * FROM people WHERE id = $1", car.Owner.ID)
			if err != nil {
				return []model.Car{}, err
			}

			for rowsOwner.Next() {
				err := rowsOwner.Scan(&car.Owner.ID, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic)
				if err != nil {
					return []model.Car{}, err
				}
			}
		}
		cars = append(cars, car)
	}
	if err = rows.Err(); err != nil {
		return []model.Car{}, err
	}

	return cars, nil
}

func UpdateCarByID(id int, updateFields map[string]string) error {
	db, err := InitDB()
	if err != nil {
		logrus.Errorf("failed to initialize to db: %s", err.Error())
		return err
	}

	query := "UPDATE car SET"
	var args []interface{}
	var index = 1

	for field, value := range updateFields {
		query += " " + field + " = $" + strconv.Itoa(index) + ","
		args = append(args, value)
		index++
	}
	query = strings.TrimSuffix(query, ",")

	query += " WHERE id = $" + strconv.Itoa(index)
	args = append(args, id)

	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCarByID(id int) error {
	db, err := InitDB()
	if err != nil {
		logrus.Errorf("failed to initialize to db: %s", err.Error())
		return err
	}

	query := "DELETE FROM car WHERE id = $1"

	var carID int
	err = db.QueryRow("SELECT id FROM car WHERE id = $1", id).Scan(&carID)
	if err == sql.ErrNoRows {
		return err
	}

	_, err = db.Exec(query, id)
	if err != nil {
		logrus.Errorf("failed to delete car from db: %s", err.Error())
		return err
	}

	return nil
}
