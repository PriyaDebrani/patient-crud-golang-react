package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type Service interface {
	createPatient(p Patient) error
	getPatients() ([]Patient, error)
	getPatient(id int) (Patient, error)
	deletePatient(id int) error
	updatePatient(p Patient) error
	addSubscriber(sub Subscriber) error
	removeSubscriber(sub Subscriber) error
}

var errSubscriberNotFound = errors.New("Subscriber not found")
var errEmptySubscriber = errors.New("Subscriber name cannot be empty")

type patientsService struct {
	repo        Repository
	subscribers []Subscriber
}

type Subscriber interface {
	getName() string
	update(Notification)
}

type Notification struct {
	Message     string    `json:"message"`
	NewPatients []Patient `json:"newPatients"`
}

func newPatientsService(repo Repository) *patientsService {
	return &patientsService{
		repo:        repo,
		subscribers: []Subscriber{},
	}
}

func (s *patientsService) createPatient(p Patient) error {
	if err := patientValidation(p); err != nil {
		return err
	}

	timeNow := time.Now()
	p.CreatedAt = timeNow
	p.UpdatedAt = timeNow
	if err := s.repo.createPatient(p); err != nil {
		return err
	}

	fmt.Println("Patient created at", p.CreatedAt)
	s.notifySubscriber(fmt.Sprintf("New patient added with id: %d", p.Id))
	return nil
}

func (s *patientsService) getPatients() ([]Patient, error) {
	return s.repo.getPatients()
}

func (s *patientsService) getPatient(id int) (Patient, error) {
	return s.repo.getPatient(id)
}

func (s *patientsService) deletePatient(id int) error {
	if err := s.repo.deletePatient(id); err != nil {
		return err
	}
	s.notifySubscriber(fmt.Sprintf("Patient removed with id: %d", id))
	log.Printf("Patient removed with Id: %d", id)
	return nil
}

func (s *patientsService) updatePatient(p Patient) error {
	if err := patientValidation(p); err != nil {
		return err
	}

	p.UpdatedAt = time.Now()
	if err := s.repo.updatePatient(p); err != nil {
		return err
	}

	fmt.Println("Patient updated at", p.UpdatedAt)
	s.notifySubscriber(fmt.Sprintf("Patient updated with id: %d", p.Id))
	log.Printf("Patient updated with Id: %d", p.Id)
	return nil
}

func (s *patientsService) addSubscriber(subscriber Subscriber) error {
	if subscriber.getName() == "" {
		return errEmptySubscriber
	}
	s.subscribers = append(s.subscribers, subscriber)
	log.Printf("subscriber added: %s", subscriber.getName())
	return nil
}

func (s *patientsService) removeSubscriber(subscriber Subscriber) error {
	for i, sub := range s.subscribers {
		if sub.getName() == subscriber.getName() {
			s.subscribers = append(s.subscribers[:i], s.subscribers[i+1:]...)
			log.Printf("Subscriber removed: %s", subscriber.getName())
			return nil
		}
	}
	return errSubscriberNotFound
}

func (s *patientsService) notifySubscriber(message string) {
	patients, err := s.getPatients()
	if err != nil {
		log.Printf("Failed to get patients: %v", err)
		return
	}

	notification := Notification{
		Message:     message,
		NewPatients: patients,
	}

	for _, sub := range s.subscribers {
		sub.update(notification)
	}
}
