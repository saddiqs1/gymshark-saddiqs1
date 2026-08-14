package packs

import (
	"reflect"
	"testing"
)

func TestPacks(t *testing.T) {
	tests := map[string]struct {
		itemsOrdered int
		result       map[int]int
	}{
		"1 item": {
			itemsOrdered: 1,
			result:       map[int]int{250: 1},
		},
		"250 items": {
			itemsOrdered: 250,
			result:       map[int]int{250: 1},
		},
		"251 items": {
			itemsOrdered: 251,
			result:       map[int]int{500: 1},
		},
		"501 items": {
			itemsOrdered: 501,
			result:       map[int]int{500: 1, 250: 1},
		},
		"1001 items": {
			itemsOrdered: 1001,
			result:       map[int]int{1000: 1, 250: 1},
		},
		"12001 items": {
			itemsOrdered: 12001,
			result:       map[int]int{5000: 2, 2000: 1, 250: 1},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got, expected := Packs(test.itemsOrdered), test.result; !reflect.DeepEqual(got, expected) {
				t.Fatalf("Packs(%d) returned %v; expected %v", test.itemsOrdered, got, expected)
			}
		})
	}

}
