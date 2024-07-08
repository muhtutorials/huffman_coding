package heap

import "sort"

//         5[0]
//        /    \
//      2[1]  6[2]
//     / \     /
// 3[3] 8[4] 1[5]
//
// parent = (i-1)/2
// left child = 2i+1
// right child = 2i+2
// last non-leaf node = length/2-1

type Interface interface {
	sort.Interface
	Push(any)
	Pop() any
}

func Init(h Interface) {
	// h = []int{5, 2, 4, 3, 8, 1}
	//        5
	//      /   \
	//     2     4
	//    / \   /
	//   3   8 1
	length := h.Len() // 6
	//    last non-leaf node
	//             ↓
	for i := length/2 - 1; i > 0; i-- {
		//      2    6
		down(h, i, length)
	}
}

// down traverses the heap from parent to children
func down(h Interface, i0, n int) bool {
	// on Init
	// h = []int{5, 2, 4, 3, 8, 1}
	// i0 = index of last non-leaf node: 4[2]
	// n = array length: 6
	// "i" is parent node's index
	i := i0 // 4[2]
	//        5
	//      /   \
	//     2     4
	//    / \   /
	//   3   8 1
	for {
		j1 := 2*i + 1 // 1[5] - index of the left child of 4[2]
		//  5 >= 6 ||  5 < 0
		if j1 >= n || j1 < 0 {
			break
		}
		// j is the smallest child
		j := j1 // 1[5] - index of the left child of 4[2]
		// j2 is the right child
		//   6 = 5 + 1; check if the right child exists and is less than the left one
		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
			j = j2 // if the above condition is true then the right child is the smallest one
		}
		// if children aren't smaller than the parent then break
		if !h.Less(j, i) {
			break
		}
		//   4[2] 1[5]
		h.Swap(i, j)
		i = j // 4[5]
	}
	return i > i0 // true
}

func Push(h Interface, x any) {
	h.Push(x)
	// h.Len()-1 is index of the last element
	up(h, h.Len()-1)
}

// up traverses the heap from child to parent
func up(h Interface, j int) {
	// h = []int{1, 2, 3, 9, 8, 6, 4}
	// j = 5[7] (index of the last element)
	//        1
	//      /   \
	//     2     3
	//    / \   / \
	//   9   8 6   4
	//  /
	// 5
	for {
		// 3 = (7 - 1) / 2
		i := (j - 1) / 2 // parent of the last pushed element: 9[3]
		// 3 == 7          5[7] 9[3]
		if i == j || !h.Less(j, i) {
			break
		}
		//   9[3] 5[7]
		h.Swap(i, j)
		j = i
	}
	// h = []int{1, 2, 3, 5, 8, 6, 4, 9}
	// j = 5[3]
	//         1
	//       /   \
	//      2     3
	//     / \   / \
	//    5   8 6   4
	//   /
	//  9
	//	for {
	//		1 = (3 - 1) / 2
	//		i := (j - 1) / 2 // parent of the last pushed element: 2[1]
	//		// 1 == 3          5[3] 2[1]
	//		if i == j || !h.Less(j, i) {
	//			break
	//		}
	//	}
}

func Pop(h Interface) any {
	// h = []int{1, 2, 4, 3, 8, 6, 5}
	//        1
	//      /   \
	//     2     4
	//    / \   / \
	//   3   8 6   5
	// last index is used here instead of length because the last element
	// will be popped and there's no need to sort it
	n := h.Len() - 1 // last index = 6
	h.Swap(0, n)
	// h = []int{5, 2, 4, 3, 8, 6, 1}
	//        5
	//      /   \
	//     2     4
	//    / \   / \
	//   3   8 6   1
	down(h, 0, n)
	// removes and returns the last element in the array (1)
	return h.Pop()
}

//func down(h Interface, i0, n int) bool {
//	h = []int{5, 2, 4, 3, 8, 6, 1}
//	i0 = index of the first node: 5[0]
//	n = last element index: 6
// "i" is parent node's index
//	i := i0 // 5[0]
//        5
//      /   \
//     2     4
//    / \   / \
//   3   8 6   1
//	for {
//		j1 := 2*i + 1 // 2[1] - index of the left child of 5[0]
//		//  1 >= 6 ||  1 < 0
//		if j1 >= n || j1 < 0 {
//			break
//		}
//		// j is the smallest child
//		j := j1 // 2[1] - index of the left child of 5[0]
//		// j2 is the right child
//		//   2 = 1 + 1; check if the right child exists and is less than the left one
//				 		  2 < 6			 4[2] 2[1]
//		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
//			j = j2 // if the above condition is true then the right child is the smallest one
//		}
//		// if children aren't smaller than the parent then break
//               2[1] 5[0]
//		if !h.Less(j, i) {
//			break
//		}
//		//   5[0] 2[1]
//		h.Swap(i, j)
//		i = j // 5[1]
//	}
//	return i > i0 // false
//}

//func down(h Interface, i0, n int) bool {
//	h = []int{2, 5, 4, 3, 8, 6, 1}
//	n = last element index: 6
// "i" is parent node's index
//	i = 5[1]
//        2
//      /   \
//     5     4
//    / \   / \
//   3   8 6   1
//	for {
//		j1 := 2*i + 1 // 3[3] - index of the left child of 5[1]
//		//  3 >= 6 ||  3 < 0
//		if j1 >= n || j1 < 0 {
//			break
//		}
//		// j is the smallest child
//		j := j1 // 3[3] - index of the left child of 5[1]
//		// j2 is the right child
//		//   4 = 3 + 1; check if the right child exists and is less than the left one
//				 		  4 < 6			 8[4] 3[3]
//		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
//			j = j2 // if the above condition is true then the right child is the smallest one
//		}
//		// if children aren't smaller than the parent then break
//               3[3] 5[1]
//		if !h.Less(j, i) {
//			break
//		}
//		//   5[1] 3[3]
//		h.Swap(i, j)
//		i = j // 5[3]
//	}
//	return i > i0 // false
//}

//func down(h Interface, i0, n int) bool {
//	h = []int{2, 3, 4, 5, 8, 6, 1}
//	n = last element index: 6
// "i" is parent node's index
//	i = 5[3]
//        2
//      /   \
//     3     4
//    / \   / \
//   5   8 6   1
//	for {
//		j1 := 2*i + 1 // nil[7] - index of the left child of 5[3]
//		//  7 >= 6 ||  7 < 0
//		if j1 >= n || j1 < 0 {
//			break
//		}
//	}
//         3 > 0
//	return i > i0 // true
//}

func Remove(h Interface, i int) any {
	// []int{1, 2, 4, 3, 8, 6, 5}
	// i = 1
	//        1
	//      /   \
	//     2     4
	//    / \   / \
	//   3   8 6   5
	n := h.Len() - 1 // last index = 6
	if n != i {
		h.Swap(i, n)
		//        1
		//      /   \
		//     5     4
		//    / \   / \
		//   3   8 6   2
		if !down(h, i, n) {
			up(h, i)
		}
	}
	return h.Pop()
}

func Fix(h Interface, i int) {
	//        3
	//      /   \
	//     5     2
	//    / \   / \
	//   7   8 6   9
	if !down(h, i, h.Len()) {
		up(h, i)
	}
}

//func down(h Interface, i0, n int) bool {
//	h = []int{3, 5, 2, 7, 8, 6, 9}
//	i0 = index of the changed node: 2[2]
//	n = length of the array: 7
// "i" is parent node's index
//	i := i0 // 2[2]
//        3
//      /   \
//     5     2
//    / \   / \
//   7   8 6   9
//	for {
//		j1 := 2*i + 1 // 6[5] - index of the left child of 2[2]
//		//  5 >= 6 ||  5 < 0
//		if j1 >= n || j1 < 0 {
//			break
//		}
//		// j is the smallest child
//		j := j1 // 6[5] - index of the left child of 2[2]
//		// j2 is the right child
//		//   6 = 5 + 1; check if the right child exists and is less than the left one
//				 		  6 < 7			 9[6] 6[5]
//		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
//			j = j2 // if the above condition is true then the right child is the smallest one
//		}
//		// if children aren't smaller than the parent then break
//               6[5] 2[2]
//		if !h.Less(j, i) {
//			break
//		}
//	}
//       2[2] 2[2]
//	return i > i0 // false
//}

//func up(h Interface, j int) {
//  h = []int{3, 5, 2, 7, 8, 6, 9}
//  j = index of the changed node: 2[2]
//        3
//      /   \
//     5     2
//    / \   / \
//   7   8 6   9
//	for {
//		0 = (2 - 1) / 2
//		i := (j - 1) / 2 // parent of the element: 3[0]
//		// 0 == 2          2[2] 3[0]
//		if i == j || !h.Less(j, i) {
//			break
//		}
//		     3[0] 2[2]
//		h.Swap(i, j)
//		j = i // 2[0]
//	}
//
//  h = []int{2, 5, 3, 7, 8, 6, 9}
//  j = index of the changed node: 2[0]
//        2
//      /   \
//     5     3
//    / \   / \
//   7   8 6   9
//	for {
//		0 = (0 - 1) / 2
//		i := (j - 1) / 2 // parent of the element: 2[0]
//		// 0 == 0          2[2] 3[0]
//		if i == j || !h.Less(j, i) {
//			break
//		}
//	}
//}
