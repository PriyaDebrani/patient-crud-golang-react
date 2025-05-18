package main

import (
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type Patient struct {
	bun.BaseModel `bun:"table:patients"`

	Id        int       `json:"id" bun:"id"`
	Name      string    `json:"name" bun:"name"`
	Address   string    `json:"address" bun:"address"`
	Disease   string    `json:"disease" bun:"disease"`
	Phone     int       `json:"phone" bun:"phone"`
	Year      int       `json:"year" bun:"year"`
	Month     int       `json:"month" bun:"month"`
	Date      int       `json:"date" bun:"date"`
	CreatedAt time.Time `json:"createdAt" bun:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bun:"updated_at"`
}

const (
	mistakeNegativeId   = "id should be positive"
	mistakeInvalidMonth = "month should be positive or less than 13"
	mistakeInvalidDate  = "date should be positive or less than 32"
	mistakeInvalidyear  = "year should be positive or negative"
	mistakeEmptyName    = "name cannot be empty"
	mistakeEmptyAddress = "address cannot be empty"
	mistakeEmptyDisease = "disease cannot be empty"
	mistakeInvalidPhone = "contact Number should be postive"
)

type ValidationError struct {
	Mistakes []string `json:"mistakes"`
}

func (e ValidationError) Error() string {
	return strings.Join(e.Mistakes, ", ")
}

func patientValidation(p Patient) error {
	var mistakes []string
	if p.Id <= 0 {
		mistakes = append(mistakes, mistakeNegativeId)
	}

	if p.Name == "" {
		mistakes = append(mistakes, mistakeEmptyName)
	}

	if p.Disease == "" {
		mistakes = append(mistakes, mistakeEmptyDisease)
	}

	if p.Phone == 0 {
		mistakes = append(mistakes, mistakeInvalidPhone)
	}

	if p.Year == 0 {
		mistakes = append(mistakes, mistakeInvalidyear)
	}

	if p.Month <= 0 || p.Month > 12 {
		mistakes = append(mistakes, mistakeInvalidMonth)
	}

	if p.Date <= 0 || p.Date > 31 {
		mistakes = append(mistakes, mistakeInvalidDate)
	}

	if p.Address == "" {
		mistakes = append(mistakes, mistakeEmptyAddress)
	}
	if len(mistakes) > 0 {
		return &ValidationError{Mistakes: mistakes}
	}
	return nil
}
