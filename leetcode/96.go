package main

import "fmt"

func main() {
	fmt.Println("Trees:", Trees(5))
}

func Trees(num int) int {

	dp := make([]int, num+1)
	dp[0] = 1
	for i := 1; i <= num; i++ {
		for j := 1; j <= i; j++ {
			dp[i] += dp[j-1] * dp[i-j]
		}
		fmt.Println("dp[", i, "]:", dp[i])
	}
	return dp[num]

}
