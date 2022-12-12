package main

func main() {
	//var nums = []int{5, 4, -1, 7, 8}
	//var nums = []int{2, -1, -3, -1, 3, 0, -1, -3, 3, 0}
	//var nums = []int{-3, -2}
}

//func maxSubArray(nums []int) int {
//
//	var ans, sum int
//	ans = nums[0]
//	sum = 0
//
//	//当前最大连续子序列和为 sum
//
//	for _, num := range nums {
//		if sum > 0 {
//			//说明 sum 对结果有增益效果，则 sum 保留并加上当前遍历数字
//			sum += num
//		} else {
//			//说明 sum 对结果无增益效果，需要舍弃，则 sum 直接更新为当前遍历数字
//			sum = num
//		}
//		if sum > ans {
//			ans = sum
//		}
//	}
//
//	return ans
//}

func countOdds(low int, high int) int {

	return (high+1)/2 - low/2
}
