package main

func splitNumberByBatch(number int, batchSize int) []int {
	var result []int
	quotient := number / batchSize
	remainder := number % batchSize
	for i := 0; i < quotient; i++ {
		result = append(result, batchSize)
	}
	if remainder > 0 {
		result = append(result, remainder)
	}
	return result
}
