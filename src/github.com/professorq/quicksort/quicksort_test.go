
package quicksort

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"
    "testing"
)

func TestSwap(t *testing.T) {
    var given, expected = []int{3, 5, 8, 12}, []int{8, 5, 3, 12}
    result := given[:]
    Swap(result, 0, 2)
    for i, value := range result {
        if value != expected[i] {
            t.Logf("Result: %v, Expected %v", result, expected)
            t.FailNow()
        }
    }
}


func TestPartitionScheme(t *testing.T) {
    var tests = []struct {
        a []int             // array to be partitioned
        leading int         // First element in array
        split int           // count of (elements < leading and leading)
    }{
    {
        a: []int{10, 34, 5, 89, 212, 1, 2, 3, 4, 13},
        leading: 10,
        split: 6,
    },
    {
        a: []int{6, 7, 8, 9},
        leading: 6,
        split: 1,
    },
    }
    for _, test := range tests {
        res_split := Partition(test.a)
        if res_split != test.split {
            t.Logf("Expected index of split: %v, Actual: %v",
                   test.split, res_split)
            t.Fail()
        }
        for i, value := range test.a {
            switch {
                case i < res_split - 1:
                    if value >= test.leading {
                        t.Fail()
                    }
                case i == res_split - 1:
                    if value != test.leading {
                        t.Fail()
                    }
                case i > res_split - 1:
                    if value <= test.leading {
                        t.Fail()
                    }
            }
        }
        if t.Failed() {
            t.Log(test.a)
        }
    }
}

func TestBasicSort(t *testing.T) {
    var tests = []struct {
            m PivotChoice
            given []int
            expected []int
            name string
        }{
        {
            m: PivotLeft,
            given: []int{2, 5, 3, 4, 1},
            expected: []int{1, 2, 3, 4, 5},
            name: "Left",
        },
        {
            m: PivotMedian,
            given: []int{2, 5, 3, 4, 1},
            expected: []int{1, 2, 3, 4, 5},
            name: "Median",
        },
        {
            m: PivotRight,
            given: []int{2, 5, 3, 4, 1},
            expected: []int{1, 2, 3, 4, 5},
            name: "Right",
        },
        {
            m: PivotRand,
            given: []int{2, 5, 3, 4, 1},
            expected: []int{1, 2, 3, 4, 5},
            name: "Random",
        },
    }
    for _, test := range tests {
        Sort(test.given, test.m)
        for i, value := range test.given {
            if value != test.expected[i] {
                t.Fail()
            }
        }
        if t.Failed() {
            t.Logf("Pivot: %v", test.name)
            t.Log("Failed to sort array. Result is %v", test.given)
        }
    }
}

func TestLeftPivot(t *testing.T) {
    ints := _LoadInts("quicksort.txt")
    count := Sort(ints, PivotLeft)
    fmt.Printf("Left: %v\n", count)
    prev := -1
    for _, value := range ints {
        if value <= prev {
            t.Fail()
        }
        prev = value
    }
}

func TestMedianPivot(t *testing.T) {
    ints := _LoadInts("quicksort.txt")
    count := Sort(ints, PivotMedian)
    fmt.Printf("Median: %v\n", count)
    prev := -1
    for _, value := range ints {
        if value <= prev {
            t.Fail()
        }
        prev = value
    }
}

func TestRightPivot(t *testing.T) {
    ints := _LoadInts("quicksort.txt")
    count := Sort(ints, PivotRight)
    fmt.Printf("Right: %v\n", count)
    prev := -1
    for _, value := range ints {
        if value <= prev {
            t.Fail()
        }
        prev = value
    }
}

func TestDietersInputs(t *testing.T) {
    var tests = []struct {
        filename string
        left int
        right int
        median int
    }{
        {
            filename: "quick10.txt",
            left: 25,
            right: 29,
            median: 21,
        },
        {
            filename: "quick100.txt",
            left: 615,
            right: 587,
            median: 518,
        },
        {
            filename: "quick1000.txt",
            left: 10297,
            right: 10184,
            median: 8921,
        },
    }
    for _, test := range tests {
        left := Sort(_LoadInts(test.filename), PivotLeft)
        right := Sort(_LoadInts(test.filename), PivotRight)
        median := Sort(_LoadInts(test.filename), PivotMedian)
        if left * median * right != test.left * test.median * test.right {
            t.Logf("\nL:\tg: %v\t\tt: %v\nR:\tg: %v\t\tt: %v\nM:\tg: %v\t\tt: %v",
                   test.left, left, test.right, right, test.median, median)
            t.Fail()
        }
    }
}

func _LoadInts(filename string) (numbers []int) {
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanWords)
    for scanner.Scan() {
        x, err := strconv.Atoi(scanner.Text())
        numbers = append(numbers, x)
        if err != nil {
            fmt.Println(err)
            break
        }
    }
    return numbers
}
