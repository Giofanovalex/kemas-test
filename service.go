package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Product merepresentasikan data produk
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// MockDatabase mensimulasikan pengambilan data DB dengan latency
func MockDatabaseFetch(id int) (*Product, error) {
	time.Sleep(time.Duration(rand.Intn(400)+100) * time.Millisecond)

	if id < 0 {
		return nil, ValidationError{Msg: fmt.Sprintf("ID Produk tidak valid: %d", id)}
	}

	if id == 999 {
		return nil, NotFoundError{Msg: fmt.Sprintf("Produk dengan ID %d tidak ditemukan", id)}
	}

	return &Product{
		ID:    id,
		Name:  fmt.Sprintf("Produk-%d", id),
		Price: id * 1000,
	}, nil
}

// FetchProductsConcurrently mengambil detail produk secara bersamaan (Task 1)
func FetchProductsConcurrently(ids []int) ([]Product, error) {
	var wg sync.WaitGroup
	productChan := make(chan Product, len(ids))
	errChan := make(chan error, len(ids))

	for _, id := range ids {
		wg.Add(1)
		go func(pID int) {
			defer wg.Done()
			
			product, err := MockDatabaseFetch(pID)
			if err != nil {
				errChan <- err
				return
			}
			productChan <- *product
		}(id)
	}

	go func() {
		wg.Wait()
		close(productChan)
		close(errChan)
	}()

	results := make([]Product, 0, len(ids))
	for p := range productChan {
		results = append(results, p)
	}

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return results, nil
}
