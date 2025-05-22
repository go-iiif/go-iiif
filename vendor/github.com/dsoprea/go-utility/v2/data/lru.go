package ridata

import (
	"errors"
	"fmt"

	"github.com/dsoprea/go-logging"
)

var (
	// ErrLruEmpty indicates that the LRU is empty..
	ErrLruEmpty = errors.New("lru is empty")
)

// LruKey is the type of an LRU key.
type LruKey interface{}

// LruItem is the interface that any item we add to the LRU must satisfy.
type LruItem interface {
	Id() LruKey
}

type lruNode struct {
	before *lruNode
	after  *lruNode
	item   LruItem
}

// String will return a string representation of the node.
func (ln *lruNode) String() string {
	var beforePhrase string
	if ln.before != nil {
		beforePhrase = fmt.Sprintf("%v", ln.before.item.Id())
	} else {
		beforePhrase = "<NULL>"
	}

	var afterPhrase string
	if ln.after != nil {
		afterPhrase = fmt.Sprintf("%v", ln.after.item.Id())
	} else {
		afterPhrase = "<NULL>"
	}

	return fmt.Sprintf("[%v] BEFORE=[%s] AFTER=[%s]", ln.item.Id(), beforePhrase, afterPhrase)
}

type lruEventFunc func(id LruKey) (err error)

// Lru establises an LRU of IDs of any type.
type Lru struct {
	top     *lruNode
	bottom  *lruNode
	lookup  map[LruKey]*lruNode
	maxSize int
	dropCb  lruEventFunc
}

// NewLru returns a new instance.
func NewLru(maxSize int) *Lru {
	return &Lru{
		lookup:  make(map[LruKey]*lruNode),
		maxSize: maxSize,
	}
}

// SetDropCb sets a callback that will be triggered whenever an item ages out
// or is manually dropped.
func (lru *Lru) SetDropCb(cb lruEventFunc) {
	lru.dropCb = cb
}

// Count returns the number of items in the LRU.
func (lru *Lru) Count() int {
	return len(lru.lookup)
}

// MaxCount returns the maximum number of items the LRU can contain.
func (lru *Lru) MaxCount() int {
	return lru.maxSize
}

// IsFull will return true if at capacity.
func (lru *Lru) IsFull() bool {
	return lru.Count() == lru.maxSize
}

// Exists will do a membership check for the given key.
func (lru *Lru) Exists(id LruKey) bool {
	_, found := lru.lookup[id]
	return found
}

// FindPosition will return the numerical position in the list. Since the LRU
// will never be very large, this call is not expensive, per se. But, it *is*
// O(n) and any call to us will compound with any loops you happen to wrap us
// into.
func (lru *Lru) FindPosition(id LruKey) int {
	node, found := lru.lookup[id]
	if found == false {
		return -1
	}

	position := 0
	for ; node.before != nil; node = node.before {
		position++
	}

	return position
}

// Get touches the cache and returns the data.
func (lru *Lru) Get(id LruKey) (found bool, item LruItem, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	if node, found := lru.lookup[id]; found == true {
		_, _, err := lru.Set(node.item)
		log.PanicIf(err)

		return true, node.item, nil
	}

	return false, nil, nil
}

// Set bumps an item to the front of the LRU. It will be added if it doesn't
// already exist. If as a result of adding an item the LRU exceeds the maximum
// size, the least recently used item will be discarded.
//
// If it was not previously in the LRU, `added` will be `true`.
func (lru *Lru) Set(item LruItem) (added bool, droppedItem LruItem, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	// TODO(dustin): !! Add tests for added/droppedItem returns.

	id := item.Id()

	node, found := lru.lookup[id]

	added = (found == false)

	if found == true {
		// It's already at the front.
		if node.before == nil {
			return added, nil, nil
		}

		// If we were at the bottom, the bottom is now whatever was upstream of
		// us.
		if lru.bottom == node {
			lru.bottom = lru.bottom.before
		}

		// Prune.
		if node.before != nil {
			node.before.after = node.after
			node.before = nil
		}

		// Insert at the front.
		node.after = lru.top

		// Point the head of the list to us.
		lru.top = node
	} else {
		node = &lruNode{
			after: lru.top,
			item:  item,
		}

		lru.lookup[id] = node

		// Point the head of the list to us.
		lru.top = node
	}

	// Update the link from the downstream node.
	if node.after != nil {
		node.after.before = node
	}

	if lru.bottom == nil {
		lru.bottom = node
	}

	if len(lru.lookup) > lru.maxSize {
		lastItemId := lru.Oldest()
		lastNode := lru.lookup[lastItemId]

		found, err := lru.Drop(lastItemId)
		log.PanicIf(err)

		if found == false {
			log.Panicf("drop of old item was ineffectual")
		}

		droppedItem = lastNode.item
	}

	return added, droppedItem, nil
}

// Drop discards the given item.
func (lru *Lru) Drop(id LruKey) (found bool, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	node, found := lru.lookup[id]
	if found == false {
		return false, nil
	}

	// Keep the `top` node up-to-date.
	if node.before == nil {
		lru.top = node.after
	}

	// Keep the `bottom` node up-to-date.
	if node.after == nil {
		lru.bottom = node.before
	}

	// Detach us from the previous node and link that node to the one after us.
	if node.before != nil {
		node.before.after = node.after
	}

	delete(lru.lookup, id)

	if lru.dropCb != nil {
		err := lru.dropCb(id)
		log.PanicIf(err)
	}

	return true, nil
}

// Newest returns the most recently used ID.
func (lru *Lru) Newest() LruKey {
	if lru.top != nil {
		return lru.top.item.Id()
	}

	return nil
}

// Oldest returns the least recently used ID.
func (lru *Lru) Oldest() LruKey {
	if lru.bottom != nil {
		return lru.bottom.item.Id()
	}

	return nil
}

// All returns a list of all IDs.
func (lru *Lru) All() []LruKey {
	collected := make([]LruKey, len(lru.lookup))
	i := 0
	for value := range lru.lookup {
		collected[i] = value
		i++
	}

	return collected
}

// PopOldest will pop the oldest entry out of the LRU and return it. It will
// return ErrLruEmpty when empty.
func (lru *Lru) PopOldest() (item LruItem, err error) {
	lk := lru.Oldest()
	if lk == nil {
		return nil, ErrLruEmpty
	}

	node := lru.lookup[lk]
	if node == nil {
		log.Panicf("something went wrong resolving the oldest item")
	}

	found, err := lru.Drop(lk)
	log.PanicIf(err)

	if found == false {
		log.Panicf("something went wrong dropping the oldest item")
	}

	return node.item, nil
}

// Dump returns a list of all IDs.
func (lru *Lru) Dump() {
	fmt.Printf("Count: (%d)\n", len(lru.lookup))
	fmt.Printf("\n")

	fmt.Printf("Top: %v\n", lru.top)
	fmt.Printf("Bottom: %v\n", lru.bottom)
	fmt.Printf("\n")

	i := 0
	for ptr := lru.top; ptr != nil; ptr = ptr.after {
		fmt.Printf("%03d: %s\n", i, ptr)
		i++
	}
}
