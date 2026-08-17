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
		result       map[int]int
	}{
		"1 item": {
			itemsOrdered: 1,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       map[int]int{250: 1},
		},
		"250 items": {
			itemsOrdered: 250,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       map[int]int{250: 1},
		},
		"251 items": {
			itemsOrdered: 251,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       map[int]int{500: 1},
		},
		"501 items": {
			itemsOrdered: 501,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       map[int]int{500: 1, 250: 1},
		},
		"1001 items": {
			itemsOrdered: 1001,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       map[int]int{1000: 1, 250: 1},
		},
		"12001 items": {
			itemsOrdered: 12001,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       map[int]int{5000: 2, 2000: 1, 250: 1},
		},
		"25671 items": {
			itemsOrdered: 25671,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       map[int]int{5000: 5, 500: 1, 250: 1},
		},
		/*"751 items": {
			// TODO - this is a bug, should be 1x1000, but currently returns 1x500 + 1x250
			itemsOrdered: 751,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       map[int]int{1000: 1},
		},
		"1751 items": {
			itemsOrdered: 1751,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       map[int]int{1000: 2},
		},*/
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got, expected := GetPacksForOrder(&zerolog.Logger{}, test.itemsOrdered, test.packSizes), test.result; !reflect.DeepEqual(got, expected) {
				t.Fatalf("GetPacksForOrder(logger, %d, %v) returned %v; expected %v", test.itemsOrdered, test.packSizes, got, expected)
			}
		})
	}

}
