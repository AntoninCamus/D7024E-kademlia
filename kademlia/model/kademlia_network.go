package model

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

//KademliaNetwork is the kademlia of the KademliaNetwork DHT on which the algorithm works
type KademliaNetwork struct {
	table    *RoutingTable
	files    map[KademliaID]file
	tableMut *sync.RWMutex
	filesMut *sync.RWMutex
}

//file is the internal representation of a file
type file struct {
	value       []byte
	refreshedAt time.Time
	fileMut     *sync.Mutex
}

//NewKademliaNetwork create a new kademlia object
func NewKademliaNetwork(me Contact) *KademliaNetwork {
	return &KademliaNetwork{
		table:    NewRoutingTable(me),
		files:    make(map[KademliaID]file),
		tableMut: &sync.RWMutex{},
		filesMut: &sync.RWMutex{},
	}
}

// LOCAL (THREAD SAFE, BASIC) FUNCTIONS :

//GetIdentity returns the contact information of the host
func (kademlia *KademliaNetwork) GetIdentity() *Contact {
	return &kademlia.table.Me
}

//RegisterContact add if possible the new *contact* to the routing table
func (kademlia *KademliaNetwork) RegisterContact(contact *Contact) {
	closestContact := kademlia.GetContacts(contact.ID,1)

	if !contact.ID.equals(kademlia.GetIdentity().ID) && (closestContact == nil || closestContact[0].ID != contact.ID) {
		log.Print("Added new contact :", contact)
		kademlia.tableMut.Lock()
		// FIXME the bucket is unlimited atm, to fix directly in it
		kademlia.table.AddContact(*contact)
		kademlia.tableMut.Unlock()
	}
}

//GetContacts returns the *number* closest contacts to the *searchedID*
func (kademlia *KademliaNetwork) GetContacts(searchedID *KademliaID, number int) []Contact {
	kademlia.tableMut.RLock()
	defer kademlia.tableMut.RUnlock()
	return kademlia.table.FindClosestContacts(searchedID, number)
}

//SaveData save the content of the file *content* under the *fileID*
func (kademlia *KademliaNetwork) SaveData(fileID *KademliaID, content []byte) error {
	kademlia.filesMut.Lock()
	kademlia.files[*fileID] = file{
		value:       content,
		refreshedAt: time.Now(),
		fileMut:     &sync.Mutex{},
	}
	kademlia.filesMut.Unlock()
	return nil
}

//GetData returns the content corresponding to the *fileID*, as well as if the file was found
func (kademlia *KademliaNetwork) GetData(fileID *KademliaID) ([]byte, bool) {
	kademlia.filesMut.RLock()
	f, exists := kademlia.files[*fileID]
	if exists {
		defer func(f file) {
			f.fileMut.Lock()
			f.refreshedAt = time.Now()
			f.fileMut.Unlock()
		}(f)
	}
	kademlia.filesMut.RUnlock()
	return f.value, exists
}

func (kademlia *KademliaNetwork) PrintFileState() string {
	var ret = "["
	for _, val := range kademlia.files {
		ret += fmt.Sprintf("%s,", val.value)
	}
	ret += "]"
	return ret
}

func (kademlia *KademliaNetwork) PrintContactState() string {
	res, _ := json.Marshal(kademlia.table)
	return string(res)
}
