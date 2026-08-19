package shippingpacks

import (
	"reflect"
	"testing"

	"github.com/rs/zerolog"
)

func TestShippingPacks(t *testing.T) {
	tests := map[string]struct {
		itemsOrdered int
		packSizes    []int
		result       []ShippingPack
	}{
		"1 item": {
			itemsOrdered: 1,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       []ShippingPack{{PackSize: 250, Count: 1}},
		},
		"250 items": {
			itemsOrdered: 250,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       []ShippingPack{{PackSize: 250, Count: 1}},
		},
		"251 items": {
			itemsOrdered: 251,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       []ShippingPack{{PackSize: 500, Count: 1}},
		},
		"501 items": {
			itemsOrdered: 501,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       []ShippingPack{{PackSize: 250, Count: 1}, {PackSize: 500, Count: 1}},
		},
		"1001 items": {
			itemsOrdered: 1001,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       []ShippingPack{{PackSize: 250, Count: 1}, {PackSize: 1000, Count: 1}},
		},
		"12001 items": {
			itemsOrdered: 1201,
			packSizes:    []int{500, 200, 100, 50, 25},
			result:       []ShippingPack{{PackSize: 25, Count: 1}, {PackSize: 200, Count: 1}, {PackSize: 500, Count: 2}},
		},
		"25671 items": {
			itemsOrdered: 25671,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       []ShippingPack{{PackSize: 250, Count: 1}, {PackSize: 500, Count: 1}, {PackSize: 5000, Count: 5}},
		},
		"751 items": {
			itemsOrdered: 751,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       []ShippingPack{{PackSize: 1000, Count: 1}},
		},
		"1751 items": {
			itemsOrdered: 1751,
			packSizes:    []int{5000, 2000, 1000, 500, 250},
			result:       []ShippingPack{{PackSize: 2000, Count: 1}},
		},
		"800 items with 600,400,300": {
			itemsOrdered: 800,
			packSizes:    []int{600, 400, 300},
			result:       []ShippingPack{{PackSize: 400, Count: 2}},
		},
		"800 items with 500,400,100": {
			itemsOrdered: 800,
			packSizes:    []int{500, 400, 100},
			result:       []ShippingPack{{PackSize: 400, Count: 2}},
		},
		"200 items with 120,80,50": {
			itemsOrdered: 200,
			packSizes:    []int{120, 80, 50},
			result:       []ShippingPack{{PackSize: 80, Count: 1}, {PackSize: 120, Count: 1}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := GetShippingPacksForOrder(&zerolog.Logger{}, test.itemsOrdered, test.packSizes)
			if err != nil {
				t.Fatalf("GetPacksForOrder(logger, %d, %v) returned unexpected error: %v", test.itemsOrdered, test.packSizes, err)
			}

			if !reflect.DeepEqual(got, test.result) {
				t.Fatalf("GetPacksForOrder(logger, %d, %v) returned %v; expected %v", test.itemsOrdered, test.packSizes, got, test.result)
			}
		})
	}

}
