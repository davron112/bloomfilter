// Package client implements an rpc client for the bloomfilter, along with Add and Check methods.
package client

import (
	"errors"
	"net/rpc"

	"github.com/davron112/bloomfilter/v2/rotate"
	rpc_bf "github.com/davron112/bloomfilter/v2/rpc"
)

// Bloomfilter rpc client type
type Bloomfilter struct {
	client *rpc.Client
}

// New creates a new bloomfilter rpc client with address
func New(address string) (*Bloomfilter, error) {
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Bloomfilter{client}, nil
}

// Add element through bloomfilter rpc client
func (b *Bloomfilter) Add(elem []byte) error {
	var addOutput rpc_bf.AddOutput
	return b.client.Call("BloomfilterRPC.Add", rpc_bf.AddInput{Elems: [][]byte{elem}}, &addOutput)
}

// AddBatch adds a set of elements through bloomfilter rpc client
func (b *Bloomfilter) AddBatch(batch [][]byte) error {
	var addOutput rpc_bf.AddOutput
	return b.client.Call("BloomfilterRPC.Add", rpc_bf.AddInput{Elems: batch}, &addOutput)
}

// Check present element through bloomfilter rpc client
func (b *Bloomfilter) Check(elem []byte) (bool, error) {
	var checkOutput rpc_bf.CheckOutput
	if err := b.client.Call("BloomfilterRPC.Check", rpc_bf.CheckInput{Elems: [][]byte{elem}}, &checkOutput); err != nil {
		return false, err
	}
	for _, v := range checkOutput.Checks {
		if !v {
			return false, nil
		}
	}
	return true, nil
}

// Union element through bloomfilter rpc client with sliding bloomfilters
func (b *Bloomfilter) Union(that interface{}) (float64, error) {
	v, ok := that.(*rotate.Bloomfilter)
	if !ok {
		return -1.0, errors.New("invalide argument to Union, expected rotate.Bloomfilter")
	}
	var unionOutput rpc_bf.UnionOutput
	if err := b.client.Call("BloomfilterRPC.Union", rpc_bf.UnionInput{BF: v}, &unionOutput); err != nil {
		return -1.0, err
	}

	return unionOutput.Capacity, nil
}

// Close bloomfilter rpc client
func (b *Bloomfilter) Close() {
	b.client.Close()
}
