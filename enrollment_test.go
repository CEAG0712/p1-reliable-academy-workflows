package enrollment

import (
    "fmt"
    "sync"
	"testing"
)



func TestEnrollmentStatusMustBeActive(t *testing.T){
	expected := "ACTIVE"
	actual :=  EntitlementStatus()

	if actual != expected {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}

func TestPaymentMustBeConfirmed(t *testing.T){
	actual := IsPaymentConfirmed("PENDING")

	if actual != false{
		t.Fatalf("expected false for PENDING payment, got %t", actual)
	}
}


func TestConfirmedPaymentReturnsTrue(t *testing.T) {
    expected := true
    actual := IsPaymentConfirmed("CONFIRMED")

    if actual != expected {
        t.Fatalf("expected %t, got %t", expected, actual)
    }
}


func TestConfirmedPaymentEvent(t *testing.T) {
    payment := PaymentEvent{
        EventID:          "EVT-105",
        TenantID:         "academy-01",
        StudentID:        "student-42",
        CourseOfferingID: "python-fall-2026",
        Status:           "CONFIRMED",
    }

    actual := IsPaymentConfirmed(payment.Status)

    if !actual {
        t.Fatal("expected confirmed payment event to be accepted")
    }
}


func TestPaymentEventWithoutTenantIsInvalid(t *testing.T) {
    payment := PaymentEvent{
        EventID:          "EVT-106",
        TenantID:         "",
        StudentID:        "student-42",
        CourseOfferingID: "python-fall-2026",
        Status:           "CONFIRMED",
    }

    actual := HasRequiredEnrollmentFields(payment)

    if actual {
        t.Fatal("expected payment without tenant ID to be invalid")
    }
}


func TestRequiredEnrollmentFields(t *testing.T) {
    tests := []struct {
        name     string
        payment  PaymentEvent
        expected bool
    }{
        {
            name: "valid payment event",
            payment: PaymentEvent{
                TenantID:         "academy-01",
                StudentID:        "student-42",
                CourseOfferingID: "python-fall-2026",
            },
            expected: true,
        },
        {
            name: "missing tenant",
            payment: PaymentEvent{
                TenantID:         "",
                StudentID:        "student-42",
                CourseOfferingID: "python-fall-2026",
            },
            expected: false,
        },
        {
            name: "missing student",
            payment: PaymentEvent{
                TenantID:         "academy-01",
                StudentID:        "",
                CourseOfferingID: "python-fall-2026",
            },
            expected: false,
        },
        {
            name: "missing course offering",
            payment: PaymentEvent{
                TenantID:         "academy-01",
                StudentID:        "student-42",
                CourseOfferingID: "",
            },
            expected: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            actual := HasRequiredEnrollmentFields(tt.payment)

            if actual != tt.expected {
                t.Fatalf(
                    "expected %t, got %t",
                    tt.expected,
                    actual,
                )
            }
        })
    }
}


func TestPendingPaymentCannotProceedToEnrollment(t *testing.T) {
    payment := PaymentEvent{
        EventID:          "EVT-200",
        TenantID:         "academy-01",
        StudentID:        "student-42",
        CourseOfferingID: "python-fall-2026",
        Status:           "PENDING",
    }

    actual := CanProceedToEnrollment(payment)

    if actual {
        t.Fatal("PENDING payment must not proceed to enrollment")
    }
}


func TestConfirmedPaymentWithRequiredFieldsCanProceed(t *testing.T) {
    payment := PaymentEvent{
        EventID:          "EVT-201",
        TenantID:         "academy-01",
        StudentID:        "student-42",
        CourseOfferingID: "python-fall-2026",
        Status:           "CONFIRMED",
    }

    actual := CanProceedToEnrollment(payment)

    if !actual {
        t.Fatal("confirmed payment with required fields should proceed")
    }
}


func TestConfirmedPaymentWithoutStudentCannotProceed(t *testing.T) {
    payment := PaymentEvent{
        EventID:          "EVT-301",
        TenantID:         "academy-01",
        StudentID:        "",
        CourseOfferingID: "python-fall-2026",
        Status:           "CONFIRMED",
    }

    actual := CanProceedToEnrollment(payment)

    if actual {
        t.Fatal("confirmed payment without student ID must not proceed")
    }
}


func TestEntitlementMapStoresCorrectStudentAccess(t *testing.T) {
    entitlements := make(map[EntitlementKey]string)

    key := EntitlementKey{
        TenantID:         "academy-01",
        StudentID:        "student-42",
        CourseOfferingID: "python-fall-2026",
    }

    entitlements[key] = "ACTIVE"

    status, exists := entitlements[key]

    if !exists {
        t.Fatal("expected entitlement to exist")
    }

    if status != "ACTIVE" {
        t.Fatalf("expected ACTIVE, got %s", status)
    }
}


func TestNewFakeEnrollmentProviderInitializesState(t *testing.T) {
    provider := NewFakeEnrollmentProvider()

    if provider == nil {
        t.Fatal("expected provider to be created")
    }

    if provider.entitlements == nil {
        t.Fatal("expected entitlement map to be initialized")
    }

    if len(provider.entitlements) != 0 {
        t.Fatal("expected new provider to have no entitlements")
    }
}


func TestFakeEnrollmentProviderGrantsEntitlement(t *testing.T) {
    provider := NewFakeEnrollmentProvider()

    key := EntitlementKey{
        TenantID:         "academy-01",
        StudentID:        "student-42",
        CourseOfferingID: "python-fall-2026",
    }

    provider.Grant(key)

    status, exists := provider.entitlements[key]

    if !exists {
        t.Fatal("expected entitlement to exist after grant")
    }

    if status != "ACTIVE" {
        t.Fatalf("expected ACTIVE, got %s", status)
    }
}


func TestFakeEnrollmentProviderVerifiesActiveEntitlement(t *testing.T) {
    provider := NewFakeEnrollmentProvider()

    key := EntitlementKey{
        TenantID:         "academy-01",
        StudentID:        "student-42",
        CourseOfferingID: "python-fall-2026",
    }

    provider.Grant(key)

    actual := provider.HasActiveEntitlement(key)

    if !actual {
        t.Fatal("expected granted entitlement to be ACTIVE")
    }
}


func TestFakeEnrollmentProviderRejectsUngrantedEntitlement(t *testing.T) {
    provider := NewFakeEnrollmentProvider()

    key := EntitlementKey{
        TenantID:         "academy-01",
        StudentID:        "student-42",
        CourseOfferingID: "python-fall-2026",
    }

    actual := provider.HasActiveEntitlement(key)

    if actual {
        t.Fatal("expected ungranted entitlement to be inactive")
    }
}


func TestFakeEnrollmentProviderIsolatesEntitlements(t *testing.T) {
    provider := NewFakeEnrollmentProvider()

    grantedKey := EntitlementKey{
        TenantID:         "academy-01",
        StudentID:        "student-42",
        CourseOfferingID: "python-fall-2026",
    }

    provider.Grant(grantedKey)

    tests := []struct {
        name     string
        key      EntitlementKey
        expected bool
    }{
        {
            name:     "exact entitlement",
            key:      grantedKey,
            expected: true,
        },
        {
            name: "different tenant",
            key: EntitlementKey{
                TenantID:         "academy-02",
                StudentID:        "student-42",
                CourseOfferingID: "python-fall-2026",
            },
            expected: false,
        },
        {
            name: "different student",
            key: EntitlementKey{
                TenantID:         "academy-01",
                StudentID:        "student-99",
                CourseOfferingID: "python-fall-2026",
            },
            expected: false,
        },
        {
            name: "different course offering",
            key: EntitlementKey{
                TenantID:         "academy-01",
                StudentID:        "student-42",
                CourseOfferingID: "go-spring-2027",
            },
            expected: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            actual := provider.HasActiveEntitlement(tt.key)

            if actual != tt.expected {
                t.Fatalf(
                    "%s: expected %t, got %t",
                    tt.name,
                    tt.expected,
                    actual,
                )
            }
        })
    }
}


func TestFakeEnrollmentProviderRepeatedGrantKeepsOneEntitlement(t *testing.T) {
    provider := NewFakeEnrollmentProvider()

    key := EntitlementKey{
        TenantID:         "academy-01",
        StudentID:        "student-42",
        CourseOfferingID: "python-fall-2026",
    }

    provider.Grant(key)
    provider.Grant(key)

    if !provider.HasActiveEntitlement(key) {
        t.Fatal("expected entitlement to remain ACTIVE")
    }

    if len(provider.entitlements) != 1 {
        t.Fatalf(
            "expected exactly 1 entitlement, got %d",
            len(provider.entitlements),
        )
    }
}


func TestFakeEnrollmentProviderCountsOnlyNewGrants(t *testing.T) {
    provider := NewFakeEnrollmentProvider()

    key := EntitlementKey{
        TenantID:         "academy-01",
        StudentID:        "student-42",
        CourseOfferingID: "python-fall-2026",
    }

    if provider.grantCount != 0 {
        t.Fatalf("expected 0 grants initially, got %d", provider.grantCount)
    }

    provider.Grant(key)

    if provider.grantCount != 1 {
        t.Fatalf("expected 1 new grant, got %d", provider.grantCount)
    }

    provider.Grant(key)

    if provider.grantCount != 1 {
        t.Fatalf("duplicate grant should not increment count, got %d", provider.grantCount)
    }
}


func TestFakeEnrollmentProviderCountsDistinctGrants(t *testing.T) {
    provider := NewFakeEnrollmentProvider()

    firstKey := EntitlementKey{
        TenantID:         "academy-01",
        StudentID:        "student-42",
        CourseOfferingID: "python-fall-2026",
    }

    secondKey := EntitlementKey{
        TenantID:         "academy-01",
        StudentID:        "student-99",
        CourseOfferingID: "python-fall-2026",
    }

    provider.Grant(firstKey)
    provider.Grant(firstKey)

    provider.Grant(secondKey)
    provider.Grant(secondKey)

    if provider.grantCount != 2 {
        t.Fatalf(
            "expected 2 newly applied grants, got %d",
            provider.grantCount,
        )
    }

    if !provider.HasActiveEntitlement(firstKey) {
        t.Fatal("expected first student to have ACTIVE access")
    }

    if !provider.HasActiveEntitlement(secondKey) {
        t.Fatal("expected second student to have ACTIVE access")
    }

    if len(provider.entitlements) != 2 {
        t.Fatalf(
            "expected 2 entitlement records, got %d",
            len(provider.entitlements),
        )
    }
}


func TestFakeEnrollmentProviderConcurrentDuplicateGrants(t *testing.T) {
    provider := NewFakeEnrollmentProvider()

    key := EntitlementKey{
        TenantID:         "academy-01",
        StudentID:        "student-42",
        CourseOfferingID: "python-fall-2026",
    }

    const workers = 10

    var wg sync.WaitGroup
    wg.Add(workers)

    for i := 0; i < workers; i++ {
        go func() {
            defer wg.Done()
            provider.Grant(key)
        }()
    }

    wg.Wait()

    if !provider.HasActiveEntitlement(key) {
        t.Fatal("expected entitlement to be ACTIVE")
    }

    if provider.grantCount != 1 {
        t.Fatalf(
            "expected exactly 1 newly applied grant, got %d",
            provider.grantCount,
        )
    }
}


func TestFakeEnrollmentProviderConcurrentDistinctGrants(t *testing.T) {
    provider := NewFakeEnrollmentProvider()

    const workers = 10

    var wg sync.WaitGroup
    wg.Add(workers)

    for i := 0; i < workers; i++ {
        studentID := fmt.Sprintf("student-%d", i)

        key := EntitlementKey{
            TenantID:         "academy-01",
            StudentID:        studentID,
            CourseOfferingID: "python-fall-2026",
        }

        go func() {
            defer wg.Done()
            provider.Grant(key)
        }()
    }

    wg.Wait()

    if provider.grantCount != workers {
        t.Fatalf(
            "expected %d newly applied grants, got %d",
            workers,
            provider.grantCount,
        )
    }

    if len(provider.entitlements) != workers {
        t.Fatalf(
            "expected %d entitlements, got %d",
            workers,
            len(provider.entitlements),
        )
    }
}














