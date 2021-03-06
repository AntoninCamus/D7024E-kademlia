package model

// ContactSorter keep only a certain number of contacts, and only keep the closest to the target
type contactSorter struct {
	target   KademliaID
	contacts []Contact
}

// NewSorter returns a new instance of a ContactSorter, with a maximum of *maxSize*
func NewSorter(target KademliaID, size int) *contactSorter {
	return &contactSorter{
		target:   target,
		contacts: make([]Contact, size),
	}
}

// InsertContact insert the *contactToInsert* into the sorter.
// If the sorter is full, keep only the closer contacts.
// Returns if the sorter state was changed or not.
func (s *contactSorter) InsertContact(contactToInsert Contact) bool {
	contactToInsert.CalcDistance(&s.target)

	// We start by checking edge cases :
	for i, c := range s.contacts {
		if c.ID == nil {
			// If one of the contact is empty, it means that there is room in the sorter
			s.contacts[i] = contactToInsert
			return true
		} else if contactToInsert.ID.equals(c.ID) {
			// If the contactToInsert is found in the sorter, we can't insert it
			return false
		}
	}

	// In the case where there is no room for any more contacts :
	// We search for the worse contact of the sorter
	worsePosition := 0
	for i := range s.contacts {
		if s.contacts[worsePosition].Less(&s.contacts[i]) {
			// If the worse contact is better than this one, this one became the worse
			worsePosition = i
		}
	}

	// If the contactToInsert is better than the worse, insert it
	if contactToInsert.Less(&s.contacts[worsePosition]) {
		s.contacts[worsePosition] = contactToInsert
		return true
	} else {
		return false
	}
}

// GetContacts return the full internal state of the sorter.
func (s *contactSorter) GetContacts() []Contact {
	newList := make([]Contact, 0)
	for _, c := range s.contacts {
		if c.ID != nil {
			newList = append(newList, c)
		}
	}
	return newList
}
