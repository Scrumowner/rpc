package main

import "fmt"

func main() {
	fmt.Println(isPalindrome(121))
}

func isPalindrome(x int) bool {
	if x < 0 {
		return false
	} else if x >= 0 && x <= 9 {
		return true
	}
	i := x
	res := 0
	for i != 0 {
		num := i % 10
		i = i / 10
		res = res*10 + num
	}
	if res == x {
		return true
	}
	return false

}
