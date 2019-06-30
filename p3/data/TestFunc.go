/*
	For testing.
	Author: Kei Fukutani
	Date  : March 31, Apr 6, 2019
*/
package data

// Test chooseClosest()
func Test_chooseClosest(list []int, selfId int, sizeClosest int) []int {
	ret, _ := chooseClosest(list, selfId, sizeClosest)
	return ret
}
