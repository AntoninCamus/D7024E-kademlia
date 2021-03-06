package model

import (
	"gotest.tools/assert"
	"testing"
)

func TestContactSorter_InsertContact_ReturnsTrueWhenEmpty(t *testing.T) {
	id := NewRandomKademliaID()
	s := NewSorter(*id, 3)

	// InsertContact should return true the 3 first times
	for i := 0; i < 3; i++ {
		res := s.InsertContact(newContact(NewRandomKademliaID(), ""))
		assert.Assert(t, res)
	}
}

func TestContactSorter_InsertContact_ReturnsFalseWhenAlreadyPresent(t *testing.T) {
	targetID := NewRandomKademliaID()
	s := NewSorter(*targetID, 3)

	c := newContact(NewRandomKademliaID(), "")
	// Inserting two time the same should return true then false
	assert.Assert(t, s.InsertContact(c))
	assert.Assert(t, !s.InsertContact(c))

	//Another case :
	c1 := newContact(KademliaIDFromString("9c7dbe89c24da341bb751281926ddc11dbb656f1"), "10.0.0.201")
	c2 := newContact(KademliaIDFromString("9c7dbe89c24da341bb751281926ddc11dbb656f1"), "10.0.0.201")
	c1.CalcDistance(NewRandomKademliaID())
	c2.CalcDistance(NewRandomKademliaID())
	assert.Assert(t, s.InsertContact(c1))
	assert.Assert(t, !s.InsertContact(c2))
}

func TestContactSorter_InsertContact_InsertOnlyCloserValues(t *testing.T) {
	targetID := NewRandomKademliaID()
	sorterSize := 5
	s := NewSorter(*targetID, sorterSize)

	for i := 0; i < 100; i++ {
		addedID := NewRandomKademliaID()
		contactsBefore := s.GetContacts()

		// Let's compute the expected result :
		added := newContact(addedID, "")
		added.CalcDistance(targetID)
		expectedResult := len(contactsBefore) < sorterSize
		for _, c := range contactsBefore {
			if expectedResult {
				break
			} else {
				c.CalcDistance(targetID)
				if added.ID.equals(c.ID) {
					expectedResult = false
					break
				} else if added.Less(&c) {
					expectedResult = true
					break
				}
			}
		}

		// Verify expected results
		assert.Equal(t, s.InsertContact(newContact(addedID, "")), expectedResult)
		if expectedResult == true {
			// Expected replaced contact position
			found := false
			for _, c := range s.GetContacts() {
				if c.ID.equals(addedID) {
					found = true
					break
				}
			}
			assert.Assert(t, found)
		}
	}
}

func TestContactSorter_ManualTest(t *testing.T) {
	id := KademliaIDFromString("ffffffffffffffffffffffffffffffffffffffff")
	s := NewSorter(*id, 3)

	c4 := newContact(KademliaIDFromString("0000000000000000000000000000000000000000"), "10.0.0.201")
	c3 := newContact(KademliaIDFromString("fffffffffffffffffffffffffffffffffffffff0"), "10.0.0.201")
	c2 := newContact(KademliaIDFromString("ffffffffffffffffffffffffffffffffffffff00"), "10.0.0.201")
	c1 := newContact(KademliaIDFromString("fffffffffffffffffffffffffffffffffffff000"), "10.0.0.201")

	assert.Assert(t, s.InsertContact(c4))
	assert.Assert(t, s.InsertContact(c3))
	assert.Assert(t, s.InsertContact(c2))
	assert.Assert(t, s.InsertContact(c1))
	assert.Assert(t, !s.InsertContact(c3))

	s.GetContacts()

	assert.Assert(t, !s.InsertContact(newContact(KademliaIDFromString("210fc7bb818639ac48a4c6afa2f1581a8b9525e2"), "")))
	s.GetContacts()
}
