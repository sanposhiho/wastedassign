package a

func p(x int) int {
	return x + 1
}

func f(param int) int {
	useOutOfIf := 1212121 // want "wasted assignment"
	ret := 0              // want "wasted assignment"
	if false {
		useOutOfIf = 10 // want "wasted assignment"
		useOutOfIf = 10 // want "reassigned, but never used afterwards"

		return 0
	} else if param == 100 {
		ret = useOutOfIf
	} else if param == 200 {
		useOutOfIf = 100 // want "wasted assignment"
		useOutOfIf = 100
		useOutOfIf = p(useOutOfIf)
		useOutOfIf += 200 // want "reassigned, but never used afterwards"
	}
	useOutOfIf = 12
	println(useOutOfIf)
	useOutOfIf = 192
	useOutOfIf += 100
	useOutOfIf += 200 // want "reassigned, but never used afterwards"
	return ret
}
