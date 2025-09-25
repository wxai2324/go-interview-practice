package main

import (
	"fmt"
)

func main() {
	// Example sorted array for testing
	arr := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}

	// Test binary search
	target := 7
	index := BinarySearch(arr, target)
	fmt.Printf("BinarySearch: %d found at index %d\n", target, index)

	// Test recursive binary search
	recursiveIndex := BinarySearchRecursive(arr, target, 0, len(arr)-1)
	fmt.Printf("BinarySearchRecursive: %d found at index %d\n", target, recursiveIndex)

	// Test find insert position
	insertTarget := 8
	insertPos := FindInsertPosition(arr, insertTarget)
	fmt.Printf("FindInsertPosition: %d should be inserted at index %d\n", insertTarget, insertPos)
}

// BinarySearch performs a standard binary search to find the target in the sorted array.
// Returns the index of the target if found, or -1 if not found.
func BinarySearch(arr []int, target int) int {
	
	start := 0
	end := len(arr) - 1
	
	for start <= end {
    mid := start + (end - start) / 2
    
    if arr[mid] == target {
	    return mid
	}else if arr[mid] < target {
	    start = mid + 1
	}else {
	    end = mid - 1
	}
}

	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, start int, end int) int {
    
    if start > end {
		return -1 // Base case: return -1 if the target is not found
	}

    
    mid := start + (end - start) / 2
    
    if arr[mid] == target {
	    return mid
	}else if arr[mid] < target {
	    return BinarySearchRecursive(arr, target, mid + 1, end)
	}else {
	    return BinarySearchRecursive(arr, target, start, mid-1)
	}
    
	return -1
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
    
    
    start := 0
	end := len(arr) - 1
	
	for start <= end {
	    
	  mid := start + (end - start) / 2
	    
      if arr[mid] == target {
	    return mid
	  }else if arr[mid] < target {
	    start = mid + 1
	  }else {
	    end = mid - 1
	  }
    }
    
	return start
}
