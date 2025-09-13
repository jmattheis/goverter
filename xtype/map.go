package xtype

import (
	"golang.org/x/tools/go/types/typeutil"
)

func NewSignatureMap[T any]() *SignatureMap[T] {
	return &SignatureMap[T]{
		table: map[hashedSignature][]entry[T]{},
	}
}

// Implements a hash map with a Signature key. Uses the hash of types.Type via typeutil.
type SignatureMap[T any] struct {
	table map[hashedSignature][]entry[T]
}

type entry[T any] struct {
	key   Signature
	value T
}

type hashedSignature struct {
	source uint32
	target uint32
}

func (m *SignatureMap[T]) At(key Signature) (T, bool) {
	for _, entry := range m.table[hash(key)] {
		if key.Identical(entry.key) {
			return entry.value, true
		}
	}
	var empty T
	return empty, false
}

func (m *SignatureMap[T]) Set(key Signature, value T) {
	hashed := hash(key)
	bucket := m.table[hashed]
	for i, entry := range bucket {
		if key.Identical(entry.key) {
			bucket[i].value = value
			return
		}
	}

	m.table[hashed] = append(bucket, entry[T]{key: key, value: value})
}

func (m *SignatureMap[T]) Values() []T {
	values := make([]T, 0, len(m.table))

	for _, bucket := range m.table {
		for _, entry := range bucket {
			values = append(values, entry.value)
		}
	}

	return values
}

func hash(t Signature) hashedSignature {
	return hashedSignature{
		source: theHasher.Hash(t.Source),
		target: theHasher.Hash(t.Target),
	}
}

var theHasher typeutil.Hasher
