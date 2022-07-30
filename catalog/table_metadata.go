// this code is from https://github.com/brunocalza/go-bustub
// there is license and copyright notice in licenses/go-bustub dir

package catalog

import (
	"github.com/ryogrid/SamehadaDB/common"
	"github.com/ryogrid/SamehadaDB/execution/executors"
	"github.com/ryogrid/SamehadaDB/execution/plans"
	"github.com/ryogrid/SamehadaDB/storage/access"
	"github.com/ryogrid/SamehadaDB/storage/index"
	"github.com/ryogrid/SamehadaDB/storage/table/schema"
)

type TableMetadata struct {
	schema *schema.Schema
	name   string
	table  *access.TableHeap
	// index data class obj of each column
	// if column has no index, respond element is nil
	indexes []index.Index
	oid     uint32
}

func NewTableMetadata(schema *schema.Schema, name string, table *access.TableHeap, oid uint32) *TableMetadata {
	ret := new(TableMetadata)
	ret.schema = schema
	ret.name = name
	ret.table = table
	ret.oid = oid

	indexes := make([]index.Index, 0)
	for idx, column_ := range schema.GetColumns() {
		if column_.HasIndex() {
			// TODO: (SDB) index bucket size is common.BucketSizeOfHashIndex (auto size extending is needed...)
			//             note: one bucket is used pages for storing index key/value pairs for a column.
			//                   one page can store 512 key/value pair
			im := index.NewIndexMetadata(column_.GetColumnName()+"_index", name, schema, []uint32{uint32(idx)})
			indexes = append(indexes, index.NewLinearProbeHashTableIndex(im, table.GetBufferPoolManager(), uint32(idx), common.BucketSizeOfHashIndex, column_.IndexHeaderPageId()))
		} else {
			indexes = append(indexes, nil)
		}
	}

	ret.indexes = indexes

	return ret
}

func (t *TableMetadata) ReconstructIndexDataOfAllCol(c *Catalog, txn *access.Transaction) {
	executionEngine := &executors.ExecutionEngine{}
	executorContext := executors.NewExecutorContext(c, t.table.GetBufferPoolManager(), txn)

	// get all tuples
	outSchema := t.schema
	seqPlan := plans.NewSeqScanPlanNode(outSchema, nil, t.OID())
	results := executionEngine.Execute(seqPlan, executorContext)

	// insert index entries correspond to each tuple and column to each index objects
	for _, index_ := range t.indexes {
		if index_ != nil {
			for _, tuple_ := range results {
				rid := tuple_.GetRID()
				index_.InsertEntry(tuple_, *rid, txn)
			}
		}
	}
}

func (t *TableMetadata) Schema() *schema.Schema {
	return t.schema
}

func (t *TableMetadata) OID() uint32 {
	return t.oid
}

func (t *TableMetadata) Table() *access.TableHeap {
	return t.table
}

func (t *TableMetadata) GetIndex(colIndex int) index.Index {
	ret := t.indexes[colIndex]
	if ret == nil {
		return nil
	} else {
		return t.indexes[colIndex]
	}
}

func (t *TableMetadata) GetColumnNum() uint32 {
	return t.schema.GetColumnCount()
}
