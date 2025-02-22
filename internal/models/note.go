package models

import "time"

type Note struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	PatientID  uint           `json:"patient_id"`
	Subjective SOAPSubjective `json:"subjective"`
	Objective  SOAPObjective  `json:"objective"`
	Assessment SOAPAssessment `json:"assessment"`
	Plan       SOAPPlan       `json:"plan"`
	VoiceData  []byte         `json:"voice_data"`
	CreatedAt  time.Time      `json:"created_at"`
}

type SOAPSubjective struct {
	PatientInfo    string   `json:"patient_info"`
	ChiefComplaint string   `json:"chief_complaint"`
	Duration       string   `json:"duration"`
	History        string   `json:"history"`
	Symptoms       []string `json:"symptoms"`
}

type SOAPObjective struct {
	VitalSigns struct {
		Temperature      string `json:"temperature"`
		GeneralCondition string `json:"general_condition"`
	} `json:"vital_signs"`
	ExaminationFindings []string `json:"examination_findings"`
}

type SOAPAssessment struct {
	PrimaryDiagnosis string   `json:"primary_diagnosis"`
	Differentials    []string `json:"differentials"`
}

type SOAPPlan struct {
	ImmediateTreatment []string `json:"immediate_treatment"`
	Medications        []string `json:"medications"`
	FollowUp           string   `json:"follow_up"`
	ClientEducation    []string `json:"client_education"`
}
