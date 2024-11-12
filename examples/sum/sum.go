package sum

// Sum computes the sum of all given integers.
//
// It takes a variable number of integers as input and returns their sum.
// Internally, it iterates over the provided slice of integers and adds each
// value to the result.
//
// Example:
//
//	sum := Sum(1, 2, 3, 4, 5)
//	fmt.Println(sum) // Output: 15
func Sum(nums ...int) (result int) {
	for _, n := range nums {
		result += n
	}
	return
}
