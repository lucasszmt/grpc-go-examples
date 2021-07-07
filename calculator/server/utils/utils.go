package utils

func PrimeNumberDecompose(num int32) (primNums []int32, err error) {
	var div int32 = 2
	for {
		if div == num {
			primNums = append(primNums, div)
			break
		}
		if num%div == 0 {
			num /= div
			primNums = append(primNums, div)
			continue
		}
		div++
	}
	return primNums, nil
}
