package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testSubscriber struct {
	name         string
	notification []Notification
}

func (s *testSubscriber) update(notification Notification) {
	s.notification = append(s.notification, notification)
}

func (s *testSubscriber) getName() string {
	return s.name
}

func assertPatientEqual(t *testing.T, expectedPatient, actualPatient Patient) {
	assert.Equal(t, expectedPatient.Id, actualPatient.Id, "Id should match")
	assert.Equal(t, expectedPatient.Name, actualPatient.Name, "Name should match")
	assert.Equal(t, expectedPatient.Address, actualPatient.Address, "Address should match")
	assert.Equal(t, expectedPatient.Disease, actualPatient.Disease, "Disease should match")
	assert.Equal(t, expectedPatient.Phone, actualPatient.Phone, "Phone should match")
	assert.Equal(t, expectedPatient.Year, actualPatient.Year, "Year should match")
	assert.Equal(t, expectedPatient.Month, actualPatient.Month, "Month should match")
	assert.Equal(t, expectedPatient.Date, actualPatient.Date, "Date should match")
}

func notificationsEqual(t *testing.T, expectedPatients, actualPatients []Notification) {
	for i := range expectedPatients {
		assert.Equal(t, expectedPatients[i].Message, actualPatients[i].Message, "expect message to match")

		for j := range expectedPatients[i].NewPatients {
			assertPatientEqual(t, expectedPatients[i].NewPatients[j], actualPatients[i].NewPatients[j])
		}
	}
}

func TestService_createPatient(t *testing.T) {
	type args struct {
		patient Patient
	}

	tests := []struct {
		name             string
		args             args
		existingPatients []Patient
		wantPatients     []Patient
		wantMistakes     []string
		wantErr          error
		wantNotification []Notification
		shouldSubscribe  bool
	}{
		{
			name: "invalid id :NEG",
			args: args{
				patient: Patient{
					Id:      0,
					Name:    "wer",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
			},
			wantMistakes: []string{mistakeNegativeId},
		},
		{
			name: "invalid name :NEG",
			args: args{
				patient: Patient{
					Id:      1,
					Name:    "",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
			},
			wantMistakes: []string{mistakeEmptyName},
		},
		{
			name: "invalid address :NEG",
			args: args{
				patient: Patient{
					Id:      1,
					Name:    "wer",
					Address: "",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
			},
			wantMistakes: []string{mistakeEmptyAddress},
		},
		{
			name: "invalid disease :NEG",
			args: args{
				patient: Patient{
					Id:      1,
					Name:    "wer",
					Address: "srt",
					Disease: "",
					Phone:   12345,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
			},
			wantMistakes: []string{mistakeEmptyDisease},
		},
		{
			name: "invalid phone :NEG",
			args: args{
				patient: Patient{
					Id:      1,
					Name:    "wer",
					Address: "srt",
					Disease: "fever",
					Phone:   0,
					Year:    2024,
					Month:   10,
					Date:    12,
				},
			},
			wantMistakes: []string{mistakeInvalidPhone},
		},
		{
			name: "invalid year :NEG",
			args: args{
				patient: Patient{
					Id:      1,
					Name:    "wer",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    0,
					Month:   12,
					Date:    12,
				},
			},
			wantMistakes: []string{mistakeInvalidyear},
		},
		{
			name: "invalid month :NEG",
			args: args{
				patient: Patient{
					Id:      1,
					Name:    "wer",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   15,
					Date:    12,
				},
			},
			wantMistakes: []string{mistakeInvalidMonth},
		},
		{
			name: "invalid date :NEG",
			args: args{
				patient: Patient{
					Id:      1,
					Name:    "wer",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   10,
					Date:    -5,
				},
			},
			wantMistakes: []string{mistakeInvalidDate},
		},
		{
			name: "multiple validation errors :NEG",
			args: args{
				patient: Patient{
					Id:      1,
					Name:    "",
					Address: "",
					Disease: "",
					Phone:   0,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
			},
			wantMistakes: []string{mistakeEmptyName, mistakeEmptyDisease, mistakeInvalidPhone, mistakeEmptyAddress},
		},
		{
			name: "patient created with subscriber :POS",
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
					Id:        1,
					Name:      "wer",
					Address:   "srt",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     12,
					Date:      12,
					CreatedAt: time.Date(2024, time.March, 24, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, time.March, 24, 10, 0, 0, 0, time.UTC),
				},
			},
			wantPatients: []Patient{
				{
					Id:        1,
					Name:      "wer",
					Address:   "srt",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     12,
					Date:      12,
					CreatedAt: time.Date(2024, time.March, 24, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, time.March, 24, 10, 0, 0, 0, time.UTC),
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
			wantMistakes: nil,
			wantNotification: []Notification{
				{
					Message: "New patient added with id: 2",
					NewPatients: []Patient{
						{
							Id:        1,
							Name:      "wer",
							Address:   "srt",
							Disease:   "fever",
							Phone:     12345,
							Year:      2024,
							Month:     12,
							Date:      12,
							CreatedAt: time.Date(2024, time.March, 24, 10, 0, 0, 0, time.UTC),
							UpdatedAt: time.Date(2024, time.March, 24, 10, 0, 0, 0, time.UTC),
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
				}},
			shouldSubscribe: true,
		},
		{
			name: "patient created without subscriber :POS",
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
					Id:        1,
					Name:      "wer",
					Address:   "srt",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     12,
					Date:      12,
					CreatedAt: time.Date(2024, time.August, 12, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, time.August, 12, 10, 0, 0, 0, time.UTC),
				},
			},
			wantPatients: []Patient{
				{
					Id:        1,
					Name:      "wer",
					Address:   "srt",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     12,
					Date:      12,
					CreatedAt: time.Date(2024, time.August, 12, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, time.August, 12, 10, 0, 0, 0, time.UTC),
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
			wantMistakes: nil,
		},
		{
			name: "duplicate id :NEG",
			args: args{
				patient: Patient{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryRepository()
			repo.patients = tt.existingPatients
			service := newPatientsService(repo)
			subscriber := &testSubscriber{name: "foo"}
			if tt.shouldSubscribe {
				service.addSubscriber(subscriber)
			}

			startTime := time.Now()
			gotErr := service.createPatient(tt.args.patient)
			endTime := time.Now()

			var validationErr *ValidationError
			if tt.wantMistakes != nil {
				assert.ErrorAs(t, gotErr, &validationErr, "expect error to match")
				assert.Equal(t, tt.wantMistakes, validationErr.Mistakes, "expect validation error to match")
			} else {
				assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			}

			for i, expectedPatient := range tt.wantPatients {
				assertPatientEqual(t, expectedPatient, repo.patients[i])
			}

			notificationsEqual(t, tt.wantNotification, subscriber.notification)
			if len(tt.wantPatients) > len(tt.existingPatients) {
				for i, gotPatient := range repo.patients {
					if i < len(tt.existingPatients) {
						wantPatient := tt.existingPatients[i]
						assert.Equal(t, wantPatient.CreatedAt, gotPatient.CreatedAt, "CreatedAt should match for existing patients")
						assert.Equal(t, wantPatient.UpdatedAt, gotPatient.UpdatedAt, "UpdatedAt should match for existing patients")
					} else {
						assert.True(t, gotPatient.CreatedAt.After(startTime), "CreatedAt should be after startTime for new patients")
						assert.True(t, gotPatient.CreatedAt.Before(endTime), "CreatedAt should be before endTime for new patients")
						assert.True(t, gotPatient.UpdatedAt.After(startTime), "UpdatedAt should be after startTime for new patients")
						assert.True(t, gotPatient.UpdatedAt.Before(endTime), "UpdatedAt should be before endTime for new patients")
					}
				}
			}
		})
	}
}

func TestService_getPatients(t *testing.T) {
	tests := []struct {
		name             string
		existingPatients []Patient
		wantPatient      []Patient
	}{
		{
			name:             "empty patient list :POS",
			existingPatients: []Patient{},
			wantPatient:      []Patient{},
		},
		{
			name: "multiple patients :POS",
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
			wantPatient: []Patient{
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryRepository()
			Service := newPatientsService(repo)
			repo.patients = tt.existingPatients

			gotPatient, gotErr := Service.getPatients()

			assert.NoError(t, gotErr, "no error expected")
			assert.Equal(t, tt.wantPatient, gotPatient, "expexted and got patients are not same")
		})
	}
}

func TestService_getPatient(t *testing.T) {
	type args struct {
		id int
	}

	tests := []struct {
		name             string
		args             args
		existingPatients []Patient
		wantPatient      Patient
		wantErr          error
	}{
		{
			name: "invalid id :NEG",
			args: args{id: 2},
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
			wantPatient: Patient{},
			wantErr:     errPatientNotFound,
		},
		{
			name: "patient exists :POS",
			args: args{id: 1},
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
			wantPatient: Patient{
				Id:      1,
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
			service := newPatientsService(repo)
			repo.patients = tt.existingPatients

			gotPatient, gotErr := service.getPatient(tt.args.id)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expected error and got error are not same")
			assert.Equal(t, tt.wantPatient, gotPatient, "expected and got patients are not same")
		})
	}
}

func TestService_updatePatient(t *testing.T) {
	type args struct {
		patient Patient
	}
	tests := []struct {
		name             string
		args             args
		existingPatients []Patient
		wantPatients     []Patient
		wantErr          error
		wantNotification []Notification
		wantMistakes     []string
		shouldSubscribe  bool
	}{
		{
			name: "invalid id :NEG",
			args: args{
				patient: Patient{
					Id:      0,
					Name:    "wer",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
			},
			existingPatients: []Patient{},
			wantPatients:     []Patient{},
			wantMistakes:     []string{mistakeNegativeId},
			wantErr: &ValidationError{
				Mistakes: []string{mistakeNegativeId},
			},
		},
		{
			name: "invalid name :NEG",
			args: args{
				patient: Patient{
					Id:      1,
					Name:    "",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
			},
			existingPatients: []Patient{},
			wantPatients:     []Patient{},
			wantMistakes:     []string{mistakeEmptyName},
			wantErr: &ValidationError{
				Mistakes: []string{mistakeEmptyName},
			},
		},
		{
			name: "invalid address :NEG",
			args: args{
				patient: Patient{
					Id:      1,
					Name:    "wer",
					Address: "",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
			},
			existingPatients: []Patient{},
			wantPatients:     []Patient{},
			wantMistakes:     []string{mistakeEmptyAddress},
			wantErr: &ValidationError{
				Mistakes: []string{mistakeEmptyAddress},
			},
		},
		{
			name: "invalid disease :NEG",
			args: args{
				patient: Patient{
					Id:      1,
					Name:    "wer",
					Address: "srt",
					Disease: "",
					Phone:   12345,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
			},
			existingPatients: []Patient{},
			wantPatients:     []Patient{},
			wantMistakes:     []string{mistakeEmptyDisease},
			wantErr: &ValidationError{
				Mistakes: []string{mistakeEmptyDisease},
			},
		},
		{
			name: "invalid phone :NEG",
			args: args{
				patient: Patient{
					Id:      1,
					Name:    "wer",
					Address: "srt",
					Disease: "fever",
					Phone:   0,
					Year:    2024,
					Month:   12,
					Date:    12,
				},
			},
			existingPatients: []Patient{},
			wantPatients:     []Patient{},
			wantMistakes:     []string{mistakeInvalidPhone},
			wantErr: &ValidationError{
				Mistakes: []string{mistakeInvalidPhone},
			},
		},
		{
			name: "invalid year :NEG",
			args: args{
				patient: Patient{
					Id:      1,
					Name:    "wer",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    0,
					Month:   12,
					Date:    13,
				},
			},
			existingPatients: []Patient{},
			wantPatients:     []Patient{},
			wantMistakes:     []string{mistakeInvalidyear},
			wantErr: &ValidationError{
				Mistakes: []string{mistakeInvalidyear},
			},
		},
		{
			name: "invalid month :NEG",
			args: args{
				patient: Patient{
					Id:      1,
					Name:    "wer",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   0,
					Date:    13,
				},
			},
			existingPatients: []Patient{},
			wantPatients:     []Patient{},
			wantMistakes:     []string{mistakeInvalidMonth},
			wantErr: &ValidationError{
				Mistakes: []string{mistakeInvalidMonth},
			},
		},
		{
			name: "invalid date :NEG",
			args: args{
				patient: Patient{
					Id:      1,
					Name:    "wer",
					Address: "srt",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   10,
					Date:    -2,
				},
			},
			existingPatients: []Patient{},
			wantPatients:     []Patient{},
			wantMistakes:     []string{mistakeInvalidDate},
			wantErr: &ValidationError{
				Mistakes: []string{mistakeInvalidDate},
			},
		},
		{
			name: "patient updated with subscriber :POS",
			args: args{
				patient: Patient{
					Id:        2,
					Name:      "priya",
					Address:   "srt",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     12,
					Date:      12,
					CreatedAt: time.Date(2024, time.August, 12, 10, 0, 0, 0, time.UTC),
				},
			},
			existingPatients: []Patient{
				{
					Id:        1,
					Name:      "wer",
					Address:   "srt",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     12,
					Date:      12,
					CreatedAt: time.Date(2024, time.March, 24, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, time.March, 24, 10, 0, 0, 0, time.UTC),
				},
				{
					Id:        2,
					Name:      "abc",
					Address:   "srt",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     12,
					Date:      12,
					CreatedAt: time.Date(2024, time.August, 12, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, time.August, 12, 10, 0, 0, 0, time.UTC),
				},
			},

			wantNotification: []Notification{
				{
					Message: "Patient updated with id: 2",
					NewPatients: []Patient{
						{
							Id:        1,
							Name:      "wer",
							Address:   "srt",
							Disease:   "fever",
							Phone:     12345,
							Year:      2024,
							Month:     12,
							Date:      12,
							CreatedAt: time.Date(2024, time.March, 24, 10, 0, 0, 0, time.UTC),
							UpdatedAt: time.Date(2024, time.March, 24, 10, 0, 0, 0, time.UTC),
						},
						{
							Id:        2,
							Name:      "priya",
							Address:   "srt",
							Disease:   "fever",
							Phone:     12345,
							Year:      2024,
							Month:     12,
							Date:      12,
							CreatedAt: time.Date(2024, time.August, 12, 10, 0, 0, 0, time.UTC),
						},
					},
				}},
			shouldSubscribe: true,
		},
		{
			name: "patient updated without subscriber :POS",
			args: args{
				patient: Patient{
					Id:        2,
					Name:      "priya",
					Address:   "srt",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     12,
					Date:      12,
					CreatedAt: time.Date(2024, time.August, 12, 10, 0, 0, 0, time.UTC),
				},
			},
			existingPatients: []Patient{
				{
					Id:        1,
					Name:      "wer",
					Address:   "srt",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     12,
					Date:      12,
					CreatedAt: time.Date(2024, time.March, 24, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, time.March, 24, 10, 0, 0, 0, time.UTC),
				},
				{
					Id:        2,
					Name:      "abc",
					Address:   "srt",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     12,
					Date:      12,
					CreatedAt: time.Date(2024, time.August, 12, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, time.August, 12, 10, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryRepository()
			repo.patients = tt.existingPatients
			service := newPatientsService(repo)
			subscriber := &testSubscriber{name: "foo"}

			if tt.shouldSubscribe {
				service.addSubscriber(subscriber)
			}

			startTime := time.Now()
			gotErr := service.updatePatient(tt.args.patient)
			endTime := time.Now()

			var validationErr *ValidationError
			if tt.wantErr == nil {
				assert.NoError(t, gotErr, "error is not expected")
			} else {
				assert.ErrorAs(t, gotErr, &validationErr, "expected error and got error are not same")
				assert.Equal(t, tt.wantMistakes, validationErr.Mistakes, "expect validation error to match")
			}
			for i, expectedPatient := range tt.wantPatients {
				assertPatientEqual(t, expectedPatient, repo.patients[i])
			}

			notificationsEqual(t, tt.wantNotification, subscriber.notification)

			for i, patient := range repo.patients {
				if patient.Id == tt.args.patient.Id {
					assert.Equal(t, patient.CreatedAt, tt.args.patient.CreatedAt, "expect createdTime to match")
					assert.True(t, patient.UpdatedAt.After(startTime), "updatedTime should be after startTime")
					assert.True(t, patient.UpdatedAt.Before(endTime), "updatedTime should be before endTime")
				} else {
					prevPatients := tt.existingPatients[i]
					assert.Equal(t, prevPatients.CreatedAt, patient.CreatedAt, "CreatedAt should match for existing patients")
					assert.Equal(t, prevPatients.UpdatedAt, patient.UpdatedAt, "UpdatedAt should match for existing patients")
				}
			}
		})
	}
}

func TestService_deletePatient(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name             string
		args             args
		existingPatients []Patient
		wantPatients     []Patient
		wantErr          error
		wantNotification []Notification
		shouldSubscribe  bool
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
			name: "patient deleted without subscriber :POS",
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
		{
			name: "patient deleted with subscriber :POS",
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
			wantNotification: []Notification{
				{
					Message: "Patient removed with id: 2",
					NewPatients: []Patient{
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
				}},
			wantErr:         nil,
			shouldSubscribe: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryRepository()
			service := newPatientsService(repo)
			repo.patients = tt.existingPatients

			subscriber := &testSubscriber{name: "foo"}

			if tt.shouldSubscribe {
				service.addSubscriber(subscriber)
			}

			gotErr := service.deletePatient(tt.args.id)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expected error and got error are not same")
			assert.Equal(t, tt.wantPatients, repo.patients, "expected and got patient mismatch")
			assert.Equal(t, tt.wantNotification, subscriber.notification, "expected notification and got notification mismatched")

		})
	}
}

func TestService_addSubscriber(t *testing.T) {
	type args struct {
		sub Subscriber
	}

	tests := []struct {
		name            string
		existingSubs    []Subscriber
		args            []args
		wantSubscribers []Subscriber
		wantErr         error
	}{
		{
			name:         "add subscriber :POS",
			existingSubs: nil,
			args: []args{{
				&testSubscriber{name: "abc"},
			}},
			wantSubscribers: []Subscriber{&testSubscriber{name: "abc"}},
		},
		{
			name:         "multiple subscribers :POS",
			existingSubs: []Subscriber{&testSubscriber{name: "abc"}},
			args: []args{
				{
					&testSubscriber{name: "xyz"},
				},
				{
					&testSubscriber{name: "mnp"},
				},
			},
			wantSubscribers: []Subscriber{
				&testSubscriber{name: "abc"},
				&testSubscriber{name: "xyz"},
				&testSubscriber{name: "mnp"},
			},
		},
		{
			name:         "empty subscriber name :NEG",
			existingSubs: nil,
			args: []args{{
				&testSubscriber{name: ""},
			}},
			wantSubscribers: nil,
			wantErr:         errEmptySubscriber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryRepository()
			service := newPatientsService(repo)
			service.subscribers = tt.existingSubs

			for _, s := range tt.args {
				err := service.addSubscriber(s.sub)
				assert.ErrorIs(t, err, tt.wantErr, "expect error to match")
			}

			assert.Equal(t, tt.wantSubscribers, service.subscribers, "expect subscribers to match")
		})
	}
}

func TestService_removeSubscriber(t *testing.T) {
	type args struct {
		sub Subscriber
	}

	tests := []struct {
		name            string
		existingSubs    []Subscriber
		args            args
		wantSubscribers []Subscriber
		wantErr         error
	}{
		{
			name:            "subscriber does not exists :NEG",
			existingSubs:    []Subscriber{&testSubscriber{name: "abc"}},
			args:            args{&testSubscriber{name: "xyz"}},
			wantSubscribers: []Subscriber{&testSubscriber{name: "abc"}},
			wantErr:         errSubscriberNotFound,
		},
		{
			name:            "remove subscriber :POS",
			existingSubs:    []Subscriber{&testSubscriber{name: "abc"}, &testSubscriber{name: "xyz"}},
			args:            args{&testSubscriber{name: "abc"}},
			wantSubscribers: []Subscriber{&testSubscriber{name: "xyz"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryRepository()
			service := newPatientsService(repo)
			service.subscribers = tt.existingSubs

			gotErr := service.removeSubscriber(tt.args.sub)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantSubscribers, service.subscribers, "subscribers not equal")
		})
	}
}
