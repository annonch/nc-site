package main //telematics

import (
	"testing"
)

func TestWrite(t *testing.T) {
	WriteDB()
}

// func TestMerge1(t *testing.T) {
// 	nums1 := []int{1, 2, 3, 0, 0, 0}
// 	m := 3
// 	nums2 := []int{2, 5, 6}
// 	n := 3

// 	want := []int{1, 2, 2, 3, 5, 6}

// 	got, err := Merge(nums1, m, nums2, n)

// 	if !equal(got, want) || err != nil {
// 		t.Fatalf("you fucked it up this time, got %q, wanted %q: Error: %q", got, want, err)
// 	}
// }
