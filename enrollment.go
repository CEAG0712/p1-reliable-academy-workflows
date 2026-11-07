package enrollment

import "sync"

func EntitlementStatus() string {
	return "ACTIVE"
}

func IsPaymentConfirmed(status string) bool {
	return status == "CONFIRMED"
}

type PaymentEvent struct {
	EventID          string
	TenantID         string
	StudentID        string
	CourseOfferingID string
	Status           string
}

func HasRequiredEnrollmentFields(payment PaymentEvent) bool {
	return payment.TenantID != "" &&
		payment.StudentID != "" &&
		payment.CourseOfferingID != ""
}

func CanProceedToEnrollment(payment PaymentEvent) bool {
	return IsPaymentConfirmed(payment.Status) &&
		HasRequiredEnrollmentFields(payment)
}

type EntitlementKey struct {
	TenantID         string
	StudentID        string
	CourseOfferingID string
}

type FakeEnrollmentProvider struct {
	mu           sync.Mutex
	entitlements map[EntitlementKey]string
	grantCount   int
}

func NewFakeEnrollmentProvider() *FakeEnrollmentProvider {
	return &FakeEnrollmentProvider{
		entitlements: make(map[EntitlementKey]string),
	}
}

func (p *FakeEnrollmentProvider) Grant(key EntitlementKey) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.entitlements[key] == "ACTIVE" {
		return
	}

	p.entitlements[key] = "ACTIVE"
	p.grantCount++
}

func (p *FakeEnrollmentProvider) HasActiveEntitlement(key EntitlementKey) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	status, exists := p.entitlements[key]
	return exists && status == "ACTIVE"
}
