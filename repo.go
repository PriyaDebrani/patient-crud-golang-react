package main

import (
	"errors"
)

var errPatientNotFound = errors.New("patient not found")
var errDuplicateId = errors.New("duplicate id")

type Repository interface {
	createPatient(p Patient) error
	getPatients() ([]Patient, error)
	getPatient(id int) (Patient, error)
	deletePatient(id int) error
	updatePatient(p Patient) error
}

type InMemoryRepository struct {
	patients []Patient
}

func newInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{patients: []Patient{}}
}

func (repo *InMemoryRepository) createPatient(p Patient) error {
	for _, patient := range repo.patients {
		if patient.Id == p.Id {
			return errDuplicateId
		}
	}
	repo.patients = append(repo.patients, p)
	return nil
}

func (repo *InMemoryRepository) getPatients() ([]Patient, error) {
	patientsCopy := make([]Patient, len(repo.patients))
	copy(patientsCopy, repo.patients)
	return patientsCopy, nil
}

func (repo *InMemoryRepository) getPatient(id int) (Patient, error) {
	for _, p := range repo.patients {
		if p.Id == id {
			return p, nil
		}
	}
	return Patient{}, errPatientNotFound
}

func (repo *InMemoryRepository) deletePatient(id int) error {
	idx, err := repo.findPatientIdx(id)
	if err != nil {
		return err
	}

	sliceLen := len(repo.patients)
	lastIndex := sliceLen - 1

	if idx != lastIndex {
		repo.patients[idx] = repo.patients[lastIndex]
	}

	repo.patients = repo.patients[:lastIndex]
	return nil
}

func (repo *InMemoryRepository) updatePatient(p Patient) error {
	idx, err := repo.findPatientIdx(p.Id)
	if err != nil {
		return err
	}

	repo.patients[idx] = p
	return nil
}

func (repo *InMemoryRepository) findPatientIdx(id int) (int, error) {
	for idx, patient := range repo.patients {
		if patient.Id == id {
			return idx, nil
		}
	}
	return -1, errPatientNotFound
}
