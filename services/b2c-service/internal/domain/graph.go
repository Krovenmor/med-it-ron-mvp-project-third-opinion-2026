package domain

import (
	"strconv"
	"time"
)

type NodeType string

const (
	NodeTypeStudy          NodeType = "study"
	NodeTypeUrgentContact  NodeType = "urgent_contact"
	NodeTypeRecommendation NodeType = "recommendation"
	NodeTypeAppointment    NodeType = "appointment"
	NodeTypeVisit          NodeType = "visit"
)

type NodeStatus string

const (
	NodeStatusRecommended NodeStatus = "recommended"
	NodeStatusBooked      NodeStatus = "booked"
	NodeStatusDone        NodeStatus = "done"
)

var modalityTitles = map[string]string{
	"CT": "Компьютерная томография",
	"DX": "Рентгенография",
	"MG": "Маммография",
}

type Graph struct {
	PatientName string
	ClinicPhone string
	Nodes       []Node
	Edges       []Edge
}

type Node struct {
	ID       string
	Type     NodeType
	Title    string
	Text     string
	Date     time.Time
	Status   NodeStatus
	Mark     string
	InClinic bool
	Bookable bool
}

type Edge struct {
	From string
	To   string
}

func BuildGraph(plan Plan, history History, clinicPhone string) Graph {
	b := graphBuilder{
		graph:   Graph{PatientName: plan.PatientName, ClinicPhone: clinicPhone},
		history: history,
		used:    map[string]bool{},
	}
	for _, c := range plan.Cases {
		b.addCase(c)
	}
	return b.graph
}

type graphBuilder struct {
	graph   Graph
	history History
	used    map[string]bool
}

func (b *graphBuilder) addCase(c PlanCase) {
	studyID := "study:" + c.ID
	b.addNode(Node{ID: studyID, Type: NodeTypeStudy, Title: modalityTitles[c.Modality], Date: c.PerformedAt})

	if c.UrgentContact {
		b.addNode(Node{
			ID:    "urgent:" + c.ID,
			Type:  NodeTypeUrgentContact,
			Title: "Врач просит срочно связаться с клиникой",
			Text:  "Позвоните в клинику: " + b.graph.ClinicPhone,
		})
		b.addEdge(studyID, "urgent:"+c.ID)
	}

	for _, rec := range c.Recommendations {
		b.addRecommendation(studyID, c.PerformedAt, rec)
	}
}

func (b *graphBuilder) addRecommendation(studyID string, studyDate time.Time, rec Recommendation) {
	node := Node{
		ID:       "recommendation:" + rec.ID,
		Type:     NodeTypeRecommendation,
		Title:    rec.ServiceName,
		Text:     rec.PatientText,
		Status:   NodeStatusRecommended,
		Mark:     rec.Mark,
		InClinic: rec.InClinic(),
	}

	var outcome *Node
	if visit, ok := b.visitFor(rec, studyDate); ok {
		node.Status = NodeStatusDone
		outcome = &Node{ID: "visit:" + rec.ID, Type: NodeTypeVisit, Title: visit.ServiceName, Date: visit.VisitedAt}
	} else if appointment, ok := b.appointmentFor(rec, studyDate); ok {
		node.Status = NodeStatusBooked
		outcome = &Node{ID: "appointment:" + appointment.ID, Type: NodeTypeAppointment, Title: appointment.ServiceName, Date: appointment.ScheduledAt}
	}
	node.Bookable = node.InClinic && node.Status == NodeStatusRecommended

	b.addNode(node)
	b.addEdge(studyID, node.ID)
	if outcome != nil {
		b.addNode(*outcome)
		b.addEdge(node.ID, outcome.ID)
	}
}

func (b *graphBuilder) visitFor(rec Recommendation, studyDate time.Time) (Visit, bool) {
	if !rec.InClinic() {
		return Visit{}, false
	}
	for i, v := range b.history.Visits {
		key := "visit:" + strconv.Itoa(i)
		if !b.used[key] && v.ServiceCode == rec.ServiceCode && v.VisitedAt.After(studyDate) {
			b.used[key] = true
			return v, true
		}
	}
	return Visit{}, false
}

func (b *graphBuilder) appointmentFor(rec Recommendation, studyDate time.Time) (Appointment, bool) {
	if a, ok := b.history.AppointmentByReferral(rec.ID); ok && !b.used["appointment:"+a.ID] {
		b.used["appointment:"+a.ID] = true
		return a, true
	}
	if !rec.InClinic() {
		return Appointment{}, false
	}
	for _, a := range b.history.Appointments {
		key := "appointment:" + a.ID
		if !b.used[key] && a.ReferralID == "" && a.ServiceCode == rec.ServiceCode && a.ScheduledAt.After(studyDate) {
			b.used[key] = true
			return a, true
		}
	}
	return Appointment{}, false
}

func (b *graphBuilder) addNode(n Node) {
	b.graph.Nodes = append(b.graph.Nodes, n)
}

func (b *graphBuilder) addEdge(from, to string) {
	b.graph.Edges = append(b.graph.Edges, Edge{From: from, To: to})
}
