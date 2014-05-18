
package quicksort

import (
    "log"
    "math/rand"
)

// Basic quicksort algorithm (Hoore circa 1961)

// Sort the array in-place
func Sort(a []int, strategy PivotChoice) (count int) {
    if len(a) <= 1 {
        return
    }
    // Swap pivot to the first position in a[left]
    strategy.Pivot(a)
    index := Partition(a)
    count += len(a) -1

    left := a[:index - 1]
    count += Sort(left, strategy)

    right := a[index:]
    count += Sort(right, strategy)
    return count
}

// Insert the first element between all elements < and > the element
func Partition(a []int) (split int) {
    // Can be done in linear time with no extra overhead (memory)
    // reduces problem size (allows divide and conquer algorithm)
    // Partitioning at the base case is equivalent to sorting

    // Single scan of array - keep track of 'read' and 'unread'
    // 'read' is split into less than and greater than pivot 'p'
    // [p,  < p   |   > p     | unread    ]
    // [0,        ^split      ^edge       ]
    split = 1
    p := a[0]
    for edge, next_value := range a {
        // edge is the right-boundary of partiotioned, read sub-array
        if next_value < p {
            // insert a[edge] into the left partition (bounded by split)
            // First item to the right of split can be moved to the edge.
            if edge > split {
                Swap(a, split, edge)
            }
            // increment split
            split += 1
        }
        edge++
    }
    // The Swap the last element in the < p section with p.
    if split > 1 {
        Swap(a, 0, split -1)
    }
    return split
}

// Chooses a pivot for the array - swaps the pivot to the first element
type PivotChoice func([]int)

// Calls the function provided by the type
func (f PivotChoice) Pivot(a []int) {
    f(a)
}

// PivotLeft does nothing (but does it well)
func PivotLeft(_ []int) {
    // No-op
}

// PivotMedian uses the median value of the first, middle, and last
// element in the given array as the pivot
func PivotMedian(a []int) {
    var median int
    n := len(a)
    if n < 3 {
        return
    }
    left, mid, right := a[0], a[(n - 1)/2], a[n-1]
    switch {
        case right > left && left > mid || mid > left && left > right:
            return
        case left > mid && mid > right || right > mid && mid > left:
            median = (n - 1)/2
        case mid > right && right > left || left > right && right > mid:
            median = n-1
        default:
            log.Print(a)
            log.Panicf("left: %v, mid %v, right %v",
                       left, mid, right)
    }
    Swap(a, 0, median)
}

// PivotRight uses the last element in the slice as the pivot
func PivotRight(a []int) {
    Swap(a, 0, len(a) - 1)
}

func PivotRand(a []int) {
    i := rand.Intn(len(a))
    Swap(a, 0, i)
}

// Swap two indeces in place in a slice
func Swap(a []int, left, right int) {
    a[left], a[right] = a[right], a[left]
}
