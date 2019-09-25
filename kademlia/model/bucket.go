package model

import (
	"container/list"
)

// bucket definition
// contains a List
type bucket struct {
	list *list.List
}

// newBucket returns a new instance of a bucket
func newBucket() *bucket {
	b := &bucket{}
	b.list = list.New()
	return b
}

// addContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) addContact(contact Contact) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

// getContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *bucket) getContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.calcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

// len return the size of the bucket
func (bucket *bucket) len() int {
	return bucket.list.Len()
}
