package hospitalmanagment

import (
	"errors"
	"sync"
	"time"
)

// Priority doubles as the queue index: 0 is served first.
type Priority int

const (
	Critical Priority = iota
	High
	Medium
	Low
	numPriorities // sentinel: always last
)

var priorityNames = [numPriorities]string{"CRITICAL", "HIGH", "MEDIUM", "LOW"}

func (p Priority) String() string { return priorityNames[p] }

type DoctorType int

const (
	ConsultDoctor DoctorType = iota
	SpecialistDoctor
)

type DoctorStatus int

const (
	Available DoctorStatus = iota
	Busy
)

// One Doctor struct for both roles. Specialty is empty for consult doctors.
type Doctor struct {
	ID             int
	Name           string
	Type           DoctorType
	Specialty      string
	Status         DoctorStatus
	CurrentPatient *Patient // nil when Available
}

type Patient struct {
	ID       int
	Name     string
	Priority Priority
}

type Record struct {
	PatientID int
	DoctorID  int // 0 = reception
	Note      string
	At        time.Time
}

// PriorityQueue: one FIFO slice per priority level.
type PriorityQueue struct {
	levels [numPriorities][]*Patient
}

func (q *PriorityQueue) Push(p *Patient) {
	q.levels[p.Priority] = append(q.levels[p.Priority], p)
}

func (q *PriorityQueue) Pop() (*Patient, bool) {
	for i := range q.levels { // Critical first
		if len(q.levels[i]) > 0 {
			p := q.levels[i][0]
			q.levels[i] = q.levels[i][1:]
			return p, true
		}
	}
	return nil, false
}

type Hospital struct {
	mu              sync.Mutex
	triage          PriorityQueue             // waiting for a consult doctor
	specialtyQueues map[string]*PriorityQueue // waiting for a specialist, one queue per specialty
	doctors         []*Doctor
	records         map[int][]Record
	nextPatientID   int
	nextDoctorID    int
}

func NewHospital() *Hospital {
	return &Hospital{
		specialtyQueues: map[string]*PriorityQueue{},
		records:         map[int][]Record{},
	}
}

func (h *Hospital) AddDoctor(name string, t DoctorType, specialty string) (*Doctor, error) {
	if name == "" || (t == SpecialistDoctor && specialty == "") {
		return nil, errors.New("invalid input")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.nextDoctorID++
	d := &Doctor{ID: h.nextDoctorID, Name: name, Type: t, Specialty: specialty}
	h.doctors = append(h.doctors, d)
	if t == SpecialistDoctor {
		h.pullNextPatient(d) // a new specialist can pick up an existing backlog
	}
	return d, nil
}

// Reception
func (h *Hospital) Register(name string, p Priority) (*Patient, error) {
	if name == "" || p < 0 || p >= numPriorities {
		return nil, errors.New("invalid input")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.nextPatientID++
	pt := &Patient{ID: h.nextPatientID, Name: name, Priority: p}
	h.triage.Push(pt)
	h.addRecord(pt.ID, 0, "registered as "+p.String())
	return pt, nil
}

// Consult doctor takes the next patient. Pop and "mark busy" happen under one lock,
// so two doctors can never get the same patient.
func (h *Hospital) CallNext(d *Doctor) (*Patient, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if d.Type != ConsultDoctor {
		return nil, errors.New("not a consult doctor")
	}
	if d.Status == Busy {
		return nil, errors.New(d.Name + " is busy")
	}
	pt, ok := h.triage.Pop()
	if !ok {
		return nil, errors.New("no patients waiting")
	}
	d.Status, d.CurrentPatient = Busy, pt
	return pt, nil
}

// Consult doctor finishes: save notes, free the doctor, send the patient to a specialist.
func (h *Hospital) FinishConsult(d *Doctor, note, specialty string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if d.Type != ConsultDoctor || d.Status != Busy {
		return errors.New(d.Name + " has no patient to finish")
	}
	if !h.specialtyExists(specialty) { // validate first; doctor keeps the patient
		return errors.New("no specialist for " + specialty)
	}
	pt := d.CurrentPatient
	h.addRecord(pt.ID, d.ID, note)
	d.Status, d.CurrentPatient = Available, nil

	// The key logic: free specialist -> start now, otherwise wait in the specialty queue.
	if s := h.freeSpecialist(specialty); s != nil {
		s.Status, s.CurrentPatient = Busy, pt
	} else {
		q, ok := h.specialtyQueues[specialty]
		if !ok {
			q = &PriorityQueue{}
			h.specialtyQueues[specialty] = q
		}
		q.Push(pt)
	}
	return nil
}

// Specialist finishes: save notes, free themselves, then pull the next waiting patient.
func (h *Hospital) FinishTreatment(s *Doctor, note string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if s.Type != SpecialistDoctor || s.Status != Busy {
		return errors.New(s.Name + " has no patient to finish")
	}
	h.addRecord(s.CurrentPatient.ID, s.ID, note)
	s.Status, s.CurrentPatient = Available, nil
	h.pullNextPatient(s)
	return nil
}

func (h *Hospital) History(patientID int) []Record {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]Record(nil), h.records[patientID]...)
}

// ---- helpers (caller holds h.mu) ----

func (h *Hospital) specialtyExists(sp string) bool {
	for _, d := range h.doctors {
		if d.Type == SpecialistDoctor && d.Specialty == sp {
			return true
		}
	}
	return false
}

func (h *Hospital) freeSpecialist(sp string) *Doctor {
	for _, d := range h.doctors {
		if d.Type == SpecialistDoctor && d.Specialty == sp && d.Status == Available {
			return d
		}
	}
	return nil
}

// pullNextPatient gives a free specialist the highest-priority patient waiting for their specialty.
func (h *Hospital) pullNextPatient(s *Doctor) {
	if q, ok := h.specialtyQueues[s.Specialty]; ok {
		if pt, ok := q.Pop(); ok {
			s.Status, s.CurrentPatient = Busy, pt
		}
	}
}

func (h *Hospital) addRecord(patientID, doctorID int, note string) {
	h.records[patientID] = append(h.records[patientID], Record{patientID, doctorID, note, time.Now()})
}
