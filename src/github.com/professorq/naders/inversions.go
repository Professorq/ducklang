package main

import (
    "bufio"
    "fmt"
    "log"
    "io"
    "os"
    "strconv"
)


func IntsFrom(r io.Reader) (numbers []int) {
    scanner := bufio.NewScanner(r)
    scanner.Split(bufio.ScanWords)
    for scanner.Scan() {
        x, error := strconv.Atoi(scanner.Text())
        numbers = append(numbers, x)
        if error != nil {
            fmt.Println(error)
            return numbers
        }
    }
    return numbers
}

func CountInversions(numbers []int) (count int, merged []int) {
    median := (len(numbers) + 1) / 2
    var left =  make([]int, 0, median)
    var right = make([]int, 0, median)
    left = numbers[:median]
    right = numbers[median:]
    if len(left) == 1 && len(right) == 1 {
        if left[0] > right[0] {
            count = 1
            merged = append(merged, right[0])
            merged = append(merged, left[0])
        } else {
            merged = append(merged, left[0])
            merged = append(merged, right[0])
        }
    } else if len(left) == 0 {
        if len(right) != 1 {
            log.Fatal("Unbalanced sides of merge")
        }
        count = 1
        merged = right
    } else if len(right) == 0 {
        if len(left) != 1 {
            log.Fatal("Unbalanced sides of merge")
        }
        count = 0
        merged = left
    } else {
        // merge !!!
        l_count, l_merged := CountInversions(left)
        r_count, r_merged := CountInversions(right)
        count = r_count + l_count
        i, j := 0, 0
        for m := 0; m < len(numbers); m++ {
            if i >= len(l_merged) {
                // Left list is exhausted
                for j < len(r_merged) {
                    // Take remainder from right list
                    merged = append(merged, r_merged[j])
                    j++
                    m++
                }
                if j + i != len(numbers) {
                    log.Fatal("i and j don't add to len(numbers)")
                }
            } else if j >= len(r_merged) {
                // Right list is exhausted
                for i < len(l_merged) {
                    // take from left list only
                    merged = append(merged, l_merged[i])
                    i++
                    m++
                }
                // Small lists cannot add to original list
                if j + i != len(numbers) {
                    log.Fatal("i and j don't add to len(numbers)")
                }
            } else {
                if l_merged[i] > r_merged[j] {
                    merged = append(merged, r_merged[j])
                    j++
                    count += len(l_merged) - i
                } else if l_merged[i] < r_merged[j] {
                    merged = append(merged, l_merged[i])
                    i++
                } else {
                    log.Fatal("Something fucked, yo")
                }
            }
        }
    }
    return count, merged
}

func main() {
    filename := os.Args[1]
    file, err := os.Open(filename)
    if err != nil {
        fmt.Println(err)
    }
    ints := IntsFrom(file)
    if err != nil {
        fmt.Println(err)
    }
    inversions, list := CountInversions(ints)
    if len(list) != len(ints) {
        log.Fatal("len(list) != len(ints)")
    }
    fmt.Println(inversions)
}

