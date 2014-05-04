package mergecount

import (
    "strings"
    "testing"
)

func TestFiveElements(self *testing.T) {
    reader := strings.NewReader("2\n1\n3\n4\n5")
    numbers := IntsFrom(reader)
    inverts, sorted := CountInversions(numbers)
    if inverts != 1 {
        self.FailNow()
    }
    if !is_sorted(sorted) {
        self.FailNow()
    }
}

func TestIcanRead(self *testing.T) {
    reader := strings.NewReader("2\n1\n4\n5")
    numbers := IntsFrom(reader)
    if len(numbers) != 4 {
        self.FailNow()
    }
}

func is_sorted(numbers []int) (sorted bool) {
    prev := -1
    for i := 0; i < len(numbers); i++ {
        if numbers[i] - prev != 1 {
            return false
        prev = numbers[i]
        }
    }
    return true
}
