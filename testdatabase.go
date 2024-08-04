package gotestdatabase

import (
	"errors"
	"fmt"
)

type TestDatabase struct {
	tables        []*DatabaseNode
	tableLookup   map[string]int
	tableIDLookup map[string]int
}

func NewTestDatabase(tables []string) *TestDatabase {
	dbTables := make([]*DatabaseNode, len(tables))
	tableLookup := make(map[string]int)
	tableIDLookup := make(map[string]int)

	for k, v := range tables {
		dbTables[k] = nil
		tableLookup[v] = k
		tableIDLookup[v] = 1
	}

	return &TestDatabase{
		tables:        dbTables,
		tableLookup:   tableLookup,
		tableIDLookup: tableIDLookup,
	}
}

func (sd *TestDatabase) ListTables() []string {
	tables := make([]string, len(sd.tableLookup))
	for k, v := range sd.tableLookup {
		tables[v] = k
	}

	return tables
}

func (sd *TestDatabase) FindOne(table string, key int) (*DatabaseRecord, error) {
	tableIdx := sd.tableLookup[table]
	node := sd.tables[tableIdx]

	record := node.findOne(key)
	if record == nil {
		return nil, errors.New("Record does not exist")
	}

	return record, nil
}

func (sd *TestDatabase) InsertOne(table string, record *DatabaseRecord) int {
	tableIdx := sd.tableLookup[table]
	node := sd.tables[tableIdx]

	key := sd.tableIDLookup[table]
	sd.tableIDLookup[table]++

	sd.tables[tableIdx] = node.insertOne(key, record)

	return key
}

func (sd *TestDatabase) Insert(table string, records []*DatabaseRecord) {
	// TODO: Write a more efficient way to insert multiple records into the database and rebalance
	for _, v := range records {
		sd.InsertOne(table, v)
	}
}

func (sd *TestDatabase) UpdateOne(table string, key int, record *DatabaseRecord) {
	tableIdx := sd.tableLookup[table]
	node := sd.tables[tableIdx]

	node.updateOne(key, record)
}

func (sd *TestDatabase) RemoveOne(table string, key int) {
	tableIdx := sd.tableLookup[table]
	node := sd.tables[tableIdx]

	sd.tables[tableIdx] = node.removeOne(key)
}

type DatabaseRecord map[string]any

type DatabaseNode struct {
	key    int
	height int

	record *DatabaseRecord

	left  *DatabaseNode
	right *DatabaseNode
}

func newDatabaseNode(key int, record *DatabaseRecord, left *DatabaseNode, right *DatabaseNode) *DatabaseNode {
	return &DatabaseNode{
		key:    key,
		height: 1,
		record: record,
		left:   left,
		right:  right,
	}
}

func (n *DatabaseNode) insertOne(key int, record *DatabaseRecord) *DatabaseNode {
	if n == nil {
		return newDatabaseNode(key, record, nil, nil)
	}

	if key < n.key {
		n.left = n.left.insertOne(key, record)
	} else if key > n.key {
		n.right = n.right.insertOne(key, record)
	}

	return n.rebalanceTree()
}

func (n *DatabaseNode) removeOne(key int) *DatabaseNode {
	if n == nil {
		return nil
	}

	if key < n.key {
		n.left = n.left.removeOne(key)
	} else if key > n.key {
		n.right = n.right.removeOne(key)
	} else {
		if n.left != nil && n.right != nil {
			rightMinNode := n.right.findSmallest()
			n.key = rightMinNode.key
			n.record = rightMinNode.record

			n.right = n.right.removeOne(rightMinNode.key)
		} else if n.left != nil {
			n = n.left
		} else if n.right != nil {
			n = n.right
		} else {
			n = nil

			return n
		}
	}
	return n.rebalanceTree()
}

// Updates the record of a node, no change to tree order
func (n *DatabaseNode) updateOne(key int, record *DatabaseRecord) *DatabaseNode {
	if n == nil {
		return nil
	}

	if key < n.key {
		n.left.updateOne(key, record)
	} else if key > n.key {
		n.right.updateOne(key, record)
	}

	update := *n.record
	for k, v := range *record {
		update[k] = v
	}

	n.record = &update

	return n
}

func (n *DatabaseNode) findOne(key int) *DatabaseRecord {
	if n == nil {
		return nil
	}

	if key < n.key {
		return n.left.findOne(key)
	} else if key > n.key {
		return n.right.findOne(key)
	}

	return n.record
}

func (n *DatabaseNode) recalculateHeight() {
	n.height = 1 + max(n.left.getHeight(), n.right.getHeight())
}

func (n *DatabaseNode) getHeight() int {
	if n == nil {
		return 0
	}

	return n.height
}

func (n *DatabaseNode) rebalanceTree() *DatabaseNode {
	if n == nil {
		return n
	}
	n.recalculateHeight()

	balanceFactor := n.left.getHeight() - n.right.getHeight()
	if balanceFactor == -2 {
		if n.right.left.getHeight() > n.right.right.getHeight() {
			n.right = n.right.rotateRight()
		}

		return n.rotateLeft()
	} else if balanceFactor == 2 {
		if n.left.right.getHeight() > n.left.left.getHeight() {
			n.left = n.left.rotateLeft()
		}

		return n.rotateRight()
	}

	return n
}

func (n *DatabaseNode) rotateLeft() *DatabaseNode {
	newRoot := n.right
	n.right = newRoot.left
	newRoot.left = n

	n.recalculateHeight()
	newRoot.recalculateHeight()

	return newRoot
}

func (n *DatabaseNode) rotateRight() *DatabaseNode {
	newRoot := n.left
	n.left = newRoot.right
	newRoot.right = n

	n.recalculateHeight()
	newRoot.recalculateHeight()

	return newRoot
}

func (n *DatabaseNode) findSmallest() *DatabaseNode {
	if n.left != nil {
		return n.left.findSmallest()
	}

	return n
}

func (n *DatabaseNode) render(frame *TreeFrame, prefix string, padding string, depth int) {
	if n == nil {
		// We have nothing more to render
		return
	}

	frame.lines = append(frame.lines, padding+prefix+fmt.Sprint(n.key))

	depth++

	styler := frame.styler
	if depth != 0 {
		padding += "   "
	}

	if n.left != n {
		fork := styler.getFork()
		if n.right == nil {
			fork = styler.getLast()
		}

		leftPrefix := fork + styler.getLeft()
		n.left.render(frame, leftPrefix, padding, depth)
	}

	if n.right != nil {
		rightPrefix := styler.getLast() + styler.getRight()
		n.right.render(frame, rightPrefix, padding, depth)
	}
}
