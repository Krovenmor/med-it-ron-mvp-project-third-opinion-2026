package domain

type Urgency string

const (
	UrgencyNormal    Urgency = "normal"
	UrgencyPlanned   Urgency = "planned"
	UrgencyPriority  Urgency = "priority"
	UrgencyEmergency Urgency = "emergency"
)

func (u Urgency) Valid() bool {
	switch u {
	case UrgencyNormal, UrgencyPlanned, UrgencyPriority, UrgencyEmergency:
		return true
	}
	return false
}

type Modality string

const (
	ModalityCT Modality = "CT"
	ModalityDX Modality = "DX"
	ModalityMG Modality = "MG"
)
