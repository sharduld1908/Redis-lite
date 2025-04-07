package kvstore

import (
	"fmt"
	"hash/fnv"
)

type Node struct {
	key   string
	value any
	next  *Node
}

type HashTable struct {
	buckets  []*Node
	capacity int
	size     int
}

const (
	initialCapacity = 128 // 128 keys can be stored initially in the hashtable
	loadFactorLimit = 0.8 // Resize when hashtable is filled up 80%
)

func NewHashTable() *HashTable {
	return &HashTable{
		buckets:  make([]*Node, initialCapacity),
		capacity: initialCapacity,
		size:     0,
	}
}

func hash(key string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(key))
	return h.Sum64()
}

func (ht *HashTable) getIndex(key string) int {
	hashCode := hash(key)
	index := int(hashCode % uint64(ht.capacity))
	return index
}

func (ht *HashTable) Insert(key string, value any) {
	if float64(ht.size)/float64(ht.capacity) >= loadFactorLimit {
		ht.resize()
	}

	idx := ht.getIndex(key)
	currentNode := ht.buckets[idx]

	for currentNode != nil {
		if currentNode.key == key {
			currentNode.value = value
			return
		}
		currentNode = currentNode.next
	}

	newNode := &Node{
		key:   key,
		value: value,
		next:  ht.buckets[idx],
	}

	ht.buckets[idx] = newNode
	ht.size++
}

func (ht *HashTable) Get(key string) (any, bool) {
	index := ht.getIndex(key)
	currentNode := ht.buckets[index]

	for currentNode != nil {
		if currentNode.key == key {
			return currentNode.value, true
		}

		currentNode = currentNode.next
	}

	return nil, false
}

func (ht *HashTable) Delete(key string) {
	index := ht.getIndex(key)
	currentNode := ht.buckets[index]
	var previousNode *Node = nil

	for currentNode != nil {
		if currentNode.key == key {
			//currentNode is head
			if previousNode == nil {
				ht.buckets[index] = currentNode.next
			} else {
				previousNode.next = currentNode.next
			}

			ht.size--
			return
		}

		previousNode = currentNode
		currentNode = currentNode.next
	}

	// Key not found.
}

func (ht *HashTable) resize() {
	newCapacity := ht.capacity * 2
	newBuckets := make([]*Node, newCapacity)

	for _, headNode := range ht.buckets {
		currentNode := headNode
		for currentNode != nil {
			nextNode := currentNode.next
			newIndex := int(hash(currentNode.key) % uint64(newCapacity))
			currentNode.next = newBuckets[newIndex]
			newBuckets[newIndex] = currentNode
			currentNode = nextNode
		}
	}

	ht.buckets = newBuckets
	ht.capacity = newCapacity
	fmt.Printf("--- Resized to capacity %d ---\n", ht.capacity)
}
