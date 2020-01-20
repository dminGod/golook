// Copyright 2016 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package bolthold

import (
	"bytes"
	"reflect"
	"sort"

	bolt "go.etcd.io/bbolt"
)

// BoltholdIndexTag is the struct tag used to define an a field as indexable for a bolthold
const BoltholdIndexTag = "boltholdIndex"

const indexBucketPrefix = "_index"

// size of iterator keys stored in memory before more are fetched
const iteratorKeyMinCacheSize = 100

// Index is a function that returns the indexable, encoded bytes of the passed in value
type Index func(name string, value interface{}) ([]byte, error)

// adds an item to the index
func (s *Store) indexAdd(storer Storer, source bucketSource, key []byte, data interface{}) error {
	indexes := storer.Indexes()
	for name, index := range indexes {
		err := s.indexUpdate(storer.Type(), name, index, source, key, data, false)
		if err != nil {
			return err
		}
	}

	return nil
}

// removes an item from the index
// be sure to pass the data from the old record, not the new one
func (s *Store) indexDelete(storer Storer, source bucketSource, key []byte, originalData interface{}) error {
	indexes := storer.Indexes()

	for name, index := range indexes {
		err := s.indexUpdate(storer.Type(), name, index, source, key, originalData, true)
		if err != nil {
			return err
		}
	}

	return nil
}

// adds or removes a specific index on an item
func (s *Store) indexUpdate(typeName, indexName string, index Index, source bucketSource, key []byte, value interface{},
	delete bool) error {
	indexKey, err := index(indexName, value)
	if indexKey == nil {
		return nil
	}

	indexValue := make(keyList, 0)

	if err != nil {
		return err
	}

	b, err := source.CreateBucketIfNotExists(indexBucketName(typeName, indexName))
	if err != nil {
		return err
	}

	iVal := b.Get(indexKey)
	if iVal != nil {
		err = s.decode(iVal, &indexValue)
		if err != nil {
			return err
		}
	}

	if delete {
		indexValue.remove(key)
	} else {
		indexValue.add(key)
	}

	if len(indexValue) == 0 {
		return b.Delete(indexKey)
	}

	iVal, err = s.encode(indexValue)
	if err != nil {
		return err
	}

	return b.Put(indexKey, iVal)
}

// IndexExists tests if an index exists for the passed in field name
func (s *Store) IndexExists(source bucketSource, typeName, indexName string) bool {
	return (source.Bucket(indexBucketName(typeName, indexName)) != nil)
}

// indexBucketName returns the name of the bolt bucket where this index is stored
func indexBucketName(typeName, indexName string) []byte {
	return []byte(indexBucketPrefix + ":" + typeName + ":" + indexName)
}

// keyList is a slice of unique, sorted keys([]byte) such as what an index points to
type keyList [][]byte

func (v *keyList) add(key []byte) {
	i := sort.Search(len(*v), func(i int) bool {
		return bytes.Compare((*v)[i], key) >= 0
	})

	if i < len(*v) && bytes.Equal((*v)[i], key) {
		// already added
		return
	}

	*v = append(*v, nil)
	copy((*v)[i+1:], (*v)[i:])
	(*v)[i] = key
}

func (v *keyList) remove(key []byte) {
	i := sort.Search(len(*v), func(i int) bool {
		return bytes.Compare((*v)[i], key) >= 0
	})

	if i < len(*v) {
		copy((*v)[i:], (*v)[i+1:])
		(*v)[len(*v)-1] = nil
		*v = (*v)[:len(*v)-1]
	}
}

func (v *keyList) in(key []byte) bool {
	i := sort.Search(len(*v), func(i int) bool {
		return bytes.Compare((*v)[i], key) >= 0
	})

	return (i < len(*v) && bytes.Equal((*v)[i], key))
}

type iterator struct {
	keyCache    [][]byte
	dataBucket  *bolt.Bucket
	indexCursor *bolt.Cursor
	nextKeys    func(bool, *bolt.Cursor) ([][]byte, error)
	prepCursor  bool
	err         error
}

func (s *Store) newIterator(source bucketSource, typeName string, query *Query) *iterator {

	iter := &iterator{
		dataBucket: source.Bucket([]byte(typeName)),
		prepCursor: true,
	}

	if iter.dataBucket == nil {
		return iter
	}

	criteria := query.fieldCriteria[query.index]

	//   Key field
	if query.index == Key && !query.badIndex {
		iter.indexCursor = source.Bucket([]byte(typeName)).Cursor()

		iter.nextKeys = func(prepCursor bool, cursor *bolt.Cursor) ([][]byte, error) {
			var nKeys [][]byte

			for len(nKeys) < iteratorKeyMinCacheSize {
				var k []byte
				if prepCursor {
					k, _ = cursor.First()
					prepCursor = false
				} else {
					k, _ = cursor.Next()
				}
				if k == nil {
					return nKeys, nil
				}

				val := reflect.New(query.dataType)
				v := iter.dataBucket.Get(k)
				err := s.decode(v, val.Interface())
				if err != nil {
					return nil, err
				}

				ok, err := matchesAllCriteria(s, criteria, k, true, val.Interface())
				if err != nil {
					return nil, err
				}

				if ok {
					nKeys = append(nKeys, k)
				}
			}
			return nKeys, nil
		}

		return iter
	}

	var iBucket *bolt.Bucket
	if !query.badIndex {
		iBucket = source.Bucket(indexBucketName(typeName, query.index))
	}

	if iBucket == nil || hasMatchFunc(criteria) {
		// bad index or matches Function on indexed field, filter through entire store
		query.badIndex = true

		iter.indexCursor = source.Bucket([]byte(typeName)).Cursor()

		iter.nextKeys = func(prepCursor bool, cursor *bolt.Cursor) ([][]byte, error) {
			var nKeys [][]byte

			for len(nKeys) < iteratorKeyMinCacheSize {
				var k []byte
				if prepCursor {
					k, _ = cursor.First()
					prepCursor = false
				} else {
					k, _ = cursor.Next()
				}
				if k == nil {
					return nKeys, nil
				}

				nKeys = append(nKeys, k)
			}
			return nKeys, nil
		}

		return iter
	}

	//   indexed field
	iter.indexCursor = iBucket.Cursor()

	iter.nextKeys = func(prepCursor bool, cursor *bolt.Cursor) ([][]byte, error) {
		var nKeys [][]byte

		for len(nKeys) < iteratorKeyMinCacheSize {
			var k, v []byte
			if prepCursor {
				k, v = cursor.First()
				prepCursor = false
			} else {
				k, v = cursor.Next()
			}
			if k == nil {
				return nKeys, nil
			}

			// no currentRow on indexes as it refers to multiple rows
			ok, err := matchesAllCriteria(s, criteria, k, true, nil)
			if err != nil {
				return nil, err
			}

			if ok {
				// append the slice of keys stored in the index
				var keys = make(keyList, 0)
				err := s.decode(v, &keys)
				if err != nil {
					return nil, err
				}

				nKeys = append(nKeys, [][]byte(keys)...)
			}

		}
		return nKeys, nil

	}

	return iter

}

// Next returns the next key value that matches the iterators criteria
// If no more kv's are available the return nil, if there is an error, they return nil
// and iterator.Error() will return the error
func (i *iterator) Next() (key []byte, value []byte) {
	if i.err != nil {
		return nil, nil
	}

	if i.dataBucket == nil {
		return nil, nil
	}

	if i.nextKeys == nil {
		return nil, nil
	}

	if len(i.keyCache) == 0 {
		newKeys, err := i.nextKeys(i.prepCursor, i.indexCursor)
		i.prepCursor = false
		if err != nil {
			i.err = err
			return nil, nil
		}

		if len(newKeys) == 0 {
			return nil, nil
		}

		i.keyCache = append(i.keyCache, newKeys...)
	}

	nextKey := i.keyCache[0]
	i.keyCache = i.keyCache[1:]

	val := i.dataBucket.Get(nextKey)

	return nextKey, val
}

// Error returns the last error, iterator.Next() will not continue if there is an error present
func (i *iterator) Error() error {
	return i.err
}
