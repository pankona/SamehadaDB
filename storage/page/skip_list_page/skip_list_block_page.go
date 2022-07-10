package skip_list_page

import (
	"github.com/ryogrid/SamehadaDB/storage/buffer"
	"github.com/ryogrid/SamehadaDB/storage/page"
	"github.com/ryogrid/SamehadaDB/types"
	"unsafe"
)

// TODO: (SDB) need to modify data layout figure
// Slotted page format:
//  ---------------------------------------------------------
//  | HEADER | ... FREE SPACE ... | ... INSERTED TUPLES ... |
//  ---------------------------------------------------------
//                                ^
//                                free space pointer
//  Header format (size in bytes):
//  ----------------------------------------------------------------------------
//  | PageId (4)| LSN (4)| PrevPageId (4)| NextPageId (4)| FreeSpacePointer(4) |
//  ----------------------------------------------------------------------------
//  ----------------------------------------------------------------
//  | TupleCount (4) | Tuple_1 offset (4) | Tuple_1 size (4) | ... |
//  ----------------------------------------------------------------

const (
	DUMMY_MAX_ENTRY = 50
)

type SkipListBlockPageOnMem struct {
	//occuppied [(BlockArraySize-1)/8 + 1]byte // 256 bits
	//readable  [(BlockArraySize-1)/8 + 1]byte // 256 bits
	//array     [BlockArraySize]SkipListPair   // 252 * 16 bits
	key   []byte
	value uint32
}

type SkipListBlockPage struct {
	page.Page
	Level       int32
	SmallestKey types.Value
	Forward     []*SkipListBlockPage //[]types.PageID
	EntryCnt    int32
	MaxEntry    int32
	Entries     []*SkipListPair
	//occuppied [(BlockArraySize-1)/8 + 1]byte // 256 bits
	//readable  [(BlockArraySize-1)/8 + 1]byte // 256 bits
	//array     [BlockArraySize]SkipListPair   // 252 * 16 bits
}

func NewSkipListBlockPage(bpm *buffer.BufferPoolManager, level int32, smallestListPair *SkipListPair) *SkipListBlockPage {
	page_ := bpm.NewPage()
	if page_ == nil {
		return nil
	}

	ret := (*SkipListBlockPage)(unsafe.Pointer(page_))
	ret.Entries[0] = smallestListPair
	ret.EntryCnt = 1
	ret.MaxEntry = DUMMY_MAX_ENTRY
	ret.Forward = make([]*SkipListBlockPage, level)

	return ret
}

// TODO: (SDB) not implemented (EntryAt)

// Gets the entry at index in this node
func (node *SkipListBlockPage) EntryAt(idx int32) *SkipListPair {
	//return page_.array[index].key
	return nil
}

// TODO: (SDB) not implemented (KeyAt)

// Gets the key at index in this node
func (node *SkipListBlockPage) KeyAt(idx int32) *types.Value {
	//return page_.array[index].key
	return nil
}

// TODO: (SDB) not implemented (ValueAt)

// Gets the value at an index in this node
func (node *SkipListBlockPage) ValueAt(idx int32) uint32 {
	//return page_.array[index].value
	return 0
}

// if not found, returns info of nearest smaller key
// binary search is used for search
func (node *SkipListBlockPage) FindEntryByKey(key *types.Value) (found bool, entry *SkipListPair, index int32) {
	leftIdx := int32(-1)
	rightIdx := node.EntryCnt
	curIdx := int32(-1)
	var curEntry *SkipListPair
	for 1 < rightIdx-leftIdx {
		curIdx := (leftIdx + rightIdx) / 2
		curEntry := node.EntryAt(curIdx)
		if curEntry.Key.CompareLessThan(*key) {
			leftIdx = curIdx
		} else {
			rightIdx = curIdx
		}
	}
	//return right
	if leftIdx == rightIdx {
		return true, curEntry, curIdx
	} else {
		if curEntry.Key.CompareLessThan(*key) {
			return false, curEntry, curIdx
		} else {
			return false, node.EntryAt(curIdx - 1), curIdx - 1
		}
	}
}

// TODO: (SDB) not implemented (Insert)

// Attempts to insert a key and value into an index in the baccess.
func (node *SkipListBlockPage) Insert(key *types.Value, value uint32) bool {
	//if page_.IsOccupied(index) {
	//	return false
	//}
	//
	//page_.array[index] = SkipListPair{key, value}
	//page_.occuppied[index/8] |= (1 << (index % 8))
	//page_.readable[index/8] |= (1 << (index % 8))
	return true
}

func (node *SkipListBlockPage) Remove(key *types.Value) {
	found, _, foundIdx := node.FindEntryByKey(key)
	if found && (node.EntryCnt == 1) {
		// when there are no enry without target entry
		// this node keep reft with no entry (but new entry can be stored)

		// delete specified entry
		node.Entries[0] = nil

		var smallestKey types.Value
		switch key.ValueType() {
		case types.Integer:
			smallestKey = types.NewInteger(0)
			smallestKey.SetInfMin()
		case types.Float:
			smallestKey = types.NewFloat(0)
			smallestKey.SetInfMin()
		case types.Varchar:
			smallestKey = types.NewVarchar("")
			smallestKey.SetInfMin()
		case types.Boolean:
			smallestKey = types.NewBoolean(false)
			smallestKey.SetInfMin()
		default:
			panic("not implemented")
		}
		node.SmallestKey = smallestKey
	} else if found {
		node.Entries = append(node.Entries[:foundIdx], node.Entries[foundIdx+1:]...)
	} else { // found == false
		// do nothing
	}
}

// TODO: (SDB) not implemented (SplitNode)

// split entries of page_ at index
// new node contains entries after entry specified by index
func (node *SkipListBlockPage) SplitNode(index uint32) {
	//if !page_.IsReadable(index) {
	//	return
	//}
	//
	//page_.readable[index/8] &= ^(1 << (index % 8))
}
