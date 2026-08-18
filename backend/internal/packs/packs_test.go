package packs

import (
	"reflect"
	"testing"

	"github.com/rs/zerolog"
)

func TestPacks(t *testing.T) {
	tests := map[string]struct {
		itemsOrdered int
		packSizes    []int
		result       []PackResult
	}{
		"1 item": {
			itemsOrdered: 1,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       []PackResult{{PackSize: 250, Count: 1}},
		},
		"250 items": {
			itemsOrdered: 250,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       []PackResult{{PackSize: 250, Count: 1}},
		},
		"251 items": {
			itemsOrdered: 251,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       []PackResult{{PackSize: 500, Count: 1}},
		},
		"501 items": {
			itemsOrdered: 501,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       []PackResult{{PackSize: 250, Count: 1}, {PackSize: 500, Count: 1}},
		},
		"1001 items": {
			itemsOrdered: 1001,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       []PackResult{{PackSize: 250, Count: 1}, {PackSize: 1000, Count: 1}},
		},
		"12001 items": {
			itemsOrdered: 1201,
			packSizes:    []int{500, 200, 100, 50, 25},
			result:       []PackResult{{PackSize: 25, Count: 1}, {PackSize: 200, Count: 1}, {PackSize: 500, Count: 2}},
		},
		"25671 items": {
			itemsOrdered: 25671,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       []PackResult{{PackSize: 250, Count: 1}, {PackSize: 500, Count: 1}, {PackSize: 5000, Count: 5}},
		},
		"751 items": {
			itemsOrdered: 751,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       []PackResult{{PackSize: 1000, Count: 1}},
		},
		"1751 items": {
			itemsOrdered: 1751,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       []PackResult{{PackSize: 2000, Count: 1}},
		},
		"800 items with different pack sizes": {
			itemsOrdered: 800,
			packSizes:    []int{600, 400, 300},
			result:       []PackResult{{PackSize: 400, Count: 2}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := GetPacksForOrder(&zerolog.Logger{}, test.itemsOrdered, test.packSizes)
			if err != nil {
				t.Fatalf("GetPacksForOrder(logger, %d, %v) returned unexpected error: %v", test.itemsOrdered, test.packSizes, err)
			}

			if !reflect.DeepEqual(got, test.result) {
				t.Fatalf("GetPacksForOrder(logger, %d, %v) returned %v; expected %v", test.itemsOrdered, test.packSizes, got, test.result)
			}
		})
	}

}
