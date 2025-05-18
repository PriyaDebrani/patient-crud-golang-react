package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

type postgresRepo struct {
	db *bun.DB
}

func newPostgresRepo(db *bun.DB) *postgresRepo {
	return &postgresRepo{db: db}
}

func (dbrepo *postgresRepo) doesPatientExist(id int) (bool, error) {
	var patient Patient
	err := dbrepo.db.NewSelect().Model(&patient).Where("id = ?", id).Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (dbrepo *postgresRepo) createPatient(p Patient) error {
	_, err := dbrepo.db.NewInsert().Model(&p).Exec(context.Background())
	pgDriverErr, ok := err.(pgdriver.Error)
	if ok {
		errCode := pgDriverErr.Field('C')
		if errCode == "23505" {
			return errDuplicateId
		} else {
			return err
		}
	}
	return nil
}
func (dbrepo *postgresRepo) getPatients() ([]Patient, error) {
	patients := make([]Patient, 0)
	err := dbrepo.db.NewSelect().Model(&patients).Scan(context.Background())
	return patients, err
}

func (dbrepo *postgresRepo) getPatient(id int) (Patient, error) {
	var patient Patient
	if err := dbrepo.db.NewSelect().Model(&patient).Where("id = ?", id).Scan(context.Background()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Patient{}, errPatientNotFound
		}
		return Patient{}, err
	}

	return patient, nil
}

func (dbrepo *postgresRepo) deletePatient(id int) error {
	exists, err := dbrepo.doesPatientExist(id)
	if err != nil {
		return err
	}
	if !exists {
		return errPatientNotFound
	}

	result, err := dbrepo.db.NewDelete().Model((*Patient)(nil)).Where("id = ?", id).Exec(context.Background())
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were deleted")
	}
	return nil
}

func (dbrepo *postgresRepo) updatePatient(p Patient) error {
	exists, err := dbrepo.doesPatientExist(p.Id)
	if err != nil {
		return err
	}
	if !exists {
		return errPatientNotFound
	}

	result, err := dbrepo.db.NewUpdate().Model(&p).Where("id = ?", p.Id).Exec(context.Background())
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows were updated")
	}
	return nil
}
