package models

import "time"

type Patient struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name"`
	Species      string    `json:"species"`
	Breed        string    `json:"breed"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	Allergies    []Allergy `json:"allergies" gorm:"many2many:patient_allergies"`
	MedicalNotes []Note    `json:"medical_notes" gorm:"foreignKey:PatientID"`
}

type Allergy struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

type PatientSummary struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name"`
	Breed          string         `json:"breed"`
	DateOfBirth    time.Time      `json:"date_of_birth"`
	KeyConditions  []string       `json:"key_conditions"`
	RecentVisits   []Visit        `json:"recent_visits"`
	Medications    []Medication   `json:"current_medications"`
	Alerts         []string       `json:"alerts"`
	PreventiveCare PreventiveCare `json:"preventive_care"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	CreatedBy      string         `json:"created_by"`
	UpdatedBy      string         `json:"updated_by"`
}

type Visit struct {
	Date   string `json:"date"`
	Reason string `json:"reason"`
}

type Medication struct {
	Name      string `json:"name"`
	Dosage    string `json:"dosage"`
	Frequency string `json:"frequency"`
}

type PreventiveCare struct {
	Vaccinations []string `json:"vaccinations"`
	NextDue      []string `json:"next_due"`
}
