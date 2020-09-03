package a

func f() {
	useOutOfIf := 0 // want "wasted assignment"
	if false {
		useOutOfIf = 10 // want "wasted assignment"

		return
	}
	useOutOfIf = 12
	println(useOutOfIf)
	return
}
