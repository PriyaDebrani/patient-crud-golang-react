package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepo_createPatient(t *testing.T) {
	type args struct {
		patient Patient
	}

	tests := []struct {
		name             string
		args             args
		existingPatients []Patient
		wantPatients     []Patient
		wantErr          error
	}{
		{
			name: "duplicate patient id :NEG",
			args: args{
				patient: Patient{
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
			},
			wantPatients: []Patient{
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
			},
			wantErr: errDuplicateId,
		},
		{
			name: "succesfully add patient :POS",
			args: args{
				patient: Patient{
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
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryRepository()
			repo.patients = tt.existingPatients

			gotErr := repo.createPatient(tt.args.patient)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect errors to match")
			assert.Equal(t, repo.patients, tt.wantPatients)
		})
	}
}

func TestRepo_getPatients(t *testing.T) {
	tests := []struct {
		name             string
		existingPatients []Patient
		wantPatients     []Patient
		wantErr          bool
	}{
		{
			name:             "empty patient list :POS",
			existingPatients: []Patient{},
			wantPatients:     []Patient{},
			wantErr:          false,
		},
		{
			name: "multiple patient list :POS",
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
			wantPatients: []Patient{
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
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryRepository()
			repo.patients = tt.existingPatients

			gotPatients, gotErr := repo.getPatients()

			assert.NoError(t, gotErr, "no error expected")
			assert.Equal(t, tt.wantPatients, gotPatients, "expect patients to match")
		})
	}
}

func TestRepo_getPatient(t *testing.T) {
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
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryRepository()
			repo.patients = tt.existingPatients

			gotPatient, gotErr := repo.getPatient(tt.args.id)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expected error and got error are not same")
			assert.Equal(t, tt.wantPatient, gotPatient, "expected patient and got patient are not same")
		})
	}
}

func TestRepo_updatePatient(t *testing.T) {
	type args struct {
		patient Patient
	}

	tests := []struct {
		name             string
		existingPatients []Patient
		args             args
		wantPatients     []Patient
		wantErr          error
	}{
		{
			name:             "empty patient :NEG",
			existingPatients: []Patient{},
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
			wantErr:      errPatientNotFound,
			wantPatients: []Patient{},
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
					Id:      2,
					Name:    "xyz",
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
					Name:    "priya",
					Address: "surat",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   2,
					Date:    22,
				},
				{
					Id:      2,
					Name:    "abc",
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
					Name:    "priya",
					Address: "surat",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   2,
					Date:    22,
				},
				{
					Id:      2,
					Name:    "xyz",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryRepository()
			repo.patients = tt.existingPatients

			gotErr := repo.updatePatient(tt.args.patient)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expected error and got error are not same")
			assert.Equal(t, tt.wantPatients, repo.patients, "expected patient and got patient are not same")
		})
	}
}

func TestRepo_deletePatient(t *testing.T) {
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
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryRepository()
			repo.patients = tt.existingPatients

			gotErr := repo.deletePatient(tt.args.id)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expected error and got error are not same")
			assert.Equal(t, tt.wantPatients, repo.patients, "expected and got patient are not same")
		})
	}
}
