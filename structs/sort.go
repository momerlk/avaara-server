package structs

type CompFn[T any] func(a , b T) int

func swap[T any](element1 *T, element2 *T) {
	var val T
	val = *element1
	*element1 = *element2
	*element2 = val
}

func divideParts[T any](elements []T, below int, upper int , cmp CompFn[T]) int {
	var center T
	center = elements[upper]
	var i int
	i = below
	var j int
	for j = below; j < upper; j++ {
		if r := cmp(elements[j] , center); r == 0 || r < 0{
			swap(&elements[i], &elements[j])
			i += 1
		}
	}
	swap(&elements[i], &elements[upper])
	return i
}

func QuickSort[T any](elements []T, below int, upper int , cmp CompFn[T]) {
	if below < upper {
		var part int
		part = divideParts[T](elements, below, upper , cmp)
		QuickSort[T](elements, below, part-1 , cmp)
		QuickSort[T](elements, part+1, upper , cmp)
	}
}