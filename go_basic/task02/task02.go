package main

// This program demonstrates how to use pointers in Go.
// It defines a function that takes a pointer to an integer,
// adds 10 to the value at that pointer, and prints the result.

// func main() {
// 	// var a int = 5
// 	// add10(&a)
// 	// fmt.Println(a) // Output: 15

// 	// arr := []int{1, 2, 3, 4, 5}
// 	// doubleSliceElements(&arr)
// 	// fmt.Println(arr) // Output: [2 4 6 8 10]
// 	arr := []int{1, 2, 3, 4, 5}
// 	doubleSliceElements2(arr)
// 	fmt.Println(arr) // Output: [2 4 6 8 10]
// }

func add10(ptr *int) {
	*ptr += 10
}

// 实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
func doubleSliceElements(ptr *[]int) {
	for i := range *ptr {
		(*ptr)[i] *= 2
	}
}

func doubleSliceElements2(ptr []int) {
	for i := range ptr {
		ptr[i] *= 2
	}
}
