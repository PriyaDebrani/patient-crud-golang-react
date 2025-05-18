package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

func setup(db *bun.DB, existingPatients []Patient) error {
	_, err := db.NewDelete().Model((*Patient)(nil)).Where("true").Exec(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete existing patients: %w", err)
	}

	if existingPatients != nil {
		_, err = db.NewInsert().Model(&existingPatients).Exec(context.Background())
		if err != nil {
			return fmt.Errorf("failed to insert new patients: %w", err)
		}
	}
	return nil
}

func TestPostgresRepo_createPatient(t *testing.T) {
	testTime := time.Date(2024, time.August, 20, 17, 0, 0, 0, time.UTC)

	tests := []struct {
		name             string
		args             Patient
		existingPatients []Patient
		wantPatients     []Patient
		wantErr          error
	}{
		{
			name: "duplicate id :NEG",
			args: Patient{
				Id:      1,
				Name:    "abc",
				Disease: "cold",
				Phone:   12345,
				Address: "surat",
				Date:    12,
				Month:   12,
				Year:    2022,
			},
			existingPatients: []Patient{
				{
					Id:        1,
					Name:      "abc",
					Disease:   "cold",
					Phone:     12345,
					Address:   "surat",
					Date:      12,
					Month:     12,
					Year:      2022,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
			},
			wantPatients: []Patient{
				{
					Id:        1,
					Name:      "abc",
					Disease:   "cold",
					Phone:     12345,
					Address:   "surat",
					Date:      12,
					Month:     12,
					Year:      2022,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
			},
			wantErr: errDuplicateId,
		},
		{
			name: "successfully add patient :POS",
			args: Patient{
				Id:        3,
				Name:      "ert",
				Address:   "amd",
				Disease:   "fever",
				Phone:     65432,
				Year:      2024,
				Month:     12,
				Date:      2,
				CreatedAt: testTime,
				UpdatedAt: testTime,
			},
			existingPatients: []Patient{
				{
					Id:        2,
					Name:      "abc",
					Disease:   "cold",
					Phone:     12345,
					Address:   "surat",
					Date:      12,
					Month:     12,
					Year:      2022,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
			},
			wantPatients: []Patient{
				{
					Id:        2,
					Name:      "abc",
					Disease:   "cold",
					Phone:     12345,
					Address:   "surat",
					Date:      12,
					Month:     12,
					Year:      2022,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
				{
					Id:        3,
					Name:      "ert",
					Address:   "amd",
					Disease:   "fever",
					Phone:     65432,
					Year:      2024,
					Month:     12,
					Date:      2,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := connectDB("postgres", "password", "localhost", "postgres", 5432)
			repo := newPostgresRepo(db)

			if err := setup(repo.db, tt.existingPatients); err != nil {
				t.Fatalf("failed to setup test: %v", err)
			}

			gotErr := repo.createPatient(tt.args)
			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			var patients []Patient
			err := repo.db.NewSelect().Model(&patients).Scan(context.Background())
			if err != nil {
				t.Fatalf("failed to retrieve patients: %v", err)
			}
			assert.Equal(t, tt.wantPatients, patients, "expect patients to match")
		})
	}
}
func TestPostgresRepo_getPatients(t *testing.T) {
	tests := []struct {
		name             string
		existingPatients []Patient
		wantPatients     []Patient
	}{
		{
			name:             "empty patient list :POS",
			existingPatients: nil,
			wantPatients:     nil,
		},
		{
			name: "get all patients :POS",
			existingPatients: []Patient{
				{
					Id:      1,
					Name:    "abc",
					Disease: "cold",
					Phone:   12345,
					Address: "surat",
					Date:    12,
					Month:   12,
					Year:    2022,
				},
				{
					Id:      2,
					Name:    "def",
					Disease: "fever",
					Phone:   54321,
					Address: "abc",
					Date:    10,
					Month:   11,
					Year:    2023,
				},
			},
			wantPatients: []Patient{
				{
					Id:      1,
					Name:    "abc",
					Disease: "cold",
					Phone:   12345,
					Address: "surat",
					Date:    12,
					Month:   12,
					Year:    2022,
				},
				{
					Id:      2,
					Name:    "def",
					Disease: "fever",
					Phone:   54321,
					Address: "abc",
					Date:    10,
					Month:   11,
					Year:    2023,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := connectDB("postgres", "password", "localhost", "postgres", 5432)
			repo := newPostgresRepo(db)

			if err := setup(repo.db, tt.existingPatients); err != nil {
				t.Fatalf("failed to setup test: %v", err)
			}

			gotPatients, err := repo.getPatients()
			if err != nil {
				t.Fatalf("no error expected but got %v", err)
			}

			assert.Equal(t, tt.wantPatients, gotPatients, "expect patients to match")
		})
	}
}

func TestPostgresRepo_getPatient(t *testing.T) {
	type args struct {
		id int
	}

	tests := []struct {
		name             string
		existingPatients []Patient
		args             args
		wantPatient      Patient
		wantErr          error
	}{
		{
			name: "patient does not exist :NEG",
			existingPatients: []Patient{
				{
					Id:      1,
					Name:    "wer",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
				{
					Id:      2,
					Name:    "wer",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
			},
			args:        args{id: 3},
			wantPatient: Patient{},
			wantErr:     errPatientNotFound,
		},
		{
			name: "patient exists :POS",
			existingPatients: []Patient{
				{
					Id:      1,
					Name:    "wer",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
				{
					Id:      2,
					Name:    "wer",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
			},
			args: args{id: 2},
			wantPatient: Patient{
				Id:      2,
				Name:    "wer",
				Address: "srt",
				Disease: "fever",
				Phone:   12345,
				Year:    2024,
				Month:   12,
				Date:    12,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := connectDB("postgres", "password", "localhost", "postgres", 5432)
			repo := newPostgresRepo(db)

			if err := setup(repo.db, tt.existingPatients); err != nil {
				t.Fatalf("failed to setup test: %v", err)
			}

			gotPatient, gotErr := repo.getPatient(tt.args.id)
			assert.ErrorIs(t, gotErr, tt.wantErr, "expected error and got error are not same")
			assert.Equal(t, tt.wantPatient, gotPatient, "expect patients to match")
		})
	}
}

func TestPostgresRepo_deletePatient(t *testing.T) {
	type args struct {
		id int
	}

	tests := []struct {
		name             string
		args             args
		existingPatients []Patient
		wantPatients     []Patient
		wantErr          error
	}{
		{
			name: "patient does not exist :NEG",
			args: args{id: 2},
			existingPatients: []Patient{
				{
					Id:      1,
					Name:    "ert",
					Address: "amd",
					Disease: "fever",
					Phone:   65432,
					Year:    2024,
					Month:   12,
					Date:    2,
				},
			},
			wantPatients: []Patient{
				{
					Id:      1,
					Name:    "ert",
					Address: "amd",
					Disease: "fever",
					Phone:   65432,
					Year:    2024,
					Month:   12,
					Date:    2,
				},
			},
			wantErr: errPatientNotFound,
		},
		{
			name: "patient deleted :POS",
			args: args{id: 2},
			existingPatients: []Patient{
				{
					Id:      1,
					Name:    "ert",
					Address: "amd",
					Disease: "fever",
					Phone:   65432,
					Year:    2024,
					Month:   12,
					Date:    2,
				},
				{
					Id:      2,
					Name:    "ert",
					Address: "amd",
					Disease: "fever",
					Phone:   65432,
					Year:    2024,
					Month:   12,
					Date:    2,
				},
			},
			wantPatients: []Patient{
				{
					Id:      1,
					Name:    "ert",
					Address: "amd",
					Disease: "fever",
					Phone:   65432,
					Year:    2024,
					Month:   12,
					Date:    2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := connectDB("postgres", "password", "localhost", "postgres", 5432)
			repo := newPostgresRepo(db)

			if err := setup(repo.db, tt.existingPatients); err != nil {
				t.Fatalf("failed to setup test: %v", err)
			}

			gotErr := repo.deletePatient(tt.args.id)
			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			var patients []Patient
			err := repo.db.NewSelect().Model(&patients).Scan(context.Background())
			if err != nil {
				t.Fatalf("failed to retrieve patients: %v", err)
			}
			assert.Equal(t, tt.wantPatients, patients, "expect patients to match")
		})
	}
}
func TestPostgresRepo_updatePatient(t *testing.T) {
	type args struct {
		patient Patient
	}

	testTime := time.Date(2024, time.August, 20, 17, 0, 0, 0, time.UTC)

	tests := []struct {
		name             string
		existingPatients []Patient
		args             args
		wantPatients     []Patient
		wantErr          error
	}{
		{
			name:             "empty patient list :NEG",
			existingPatients: nil,
			args: args{
				patient: Patient{
					Id:      2,
					Name:    "wer",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
			},
			wantErr: errPatientNotFound,
		},
		{
			name: "patient not found :NEG",
			args: args{
				patient: Patient{
					Id:      2,
					Name:    "wer",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
			},
			existingPatients: []Patient{
				{
					Id:      1,
					Name:    "jhdfe",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   2,
					Date:    22,
				},
			},
			wantPatients: []Patient{
				{
					Id:      1,
					Name:    "jhdfe",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   2,
					Date:    22,
				},
			},
			wantErr: errPatientNotFound,
		},
		{
			name: "patient updated :POS",
			args: args{
				patient: Patient{
					Id:        2,
					Name:      "xyz",
					Address:   "srt",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     12,
					Date:      12,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
			},
			existingPatients: []Patient{
				{
					Id:        1,
					Name:      "priya",
					Address:   "surat",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     2,
					Date:      22,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
				{
					Id:        2,
					Name:      "abc",
					Address:   "srt",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     2,
					Date:      22,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
			},
			wantPatients: []Patient{
				{
					Id:        1,
					Name:      "priya",
					Address:   "surat",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     2,
					Date:      22,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
				{
					Id:        2,
					Name:      "xyz",
					Address:   "srt",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     12,
					Date:      12,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := connectDB("postgres", "password", "localhost", "postgres", 5432)
			repo := newPostgresRepo(db)

			if err := setup(repo.db, tt.existingPatients); err != nil {
				t.Fatalf("failed to setup test: %v", err)
			}

			gotErr := repo.updatePatient(tt.args.patient)
			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")

			var patients []Patient
			err := repo.db.NewSelect().Model(&patients).Scan(context.Background())
			if err != nil {
				t.Fatalf("failed to retrieve patients: %v", err)
			}
			assert.Equal(t, tt.wantPatients, patients, "expect patients to match")
		})
	}
}
