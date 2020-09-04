package a

func f() {
	useOutOfIf := 0 // want "wasted assignment"
	if false {
		useOutOfIf = 10 // want "wasted assignment"
		useOutOfIf = 10 // want "reassigned, but never used afterwards"

		return
	}
	useOutOfIf = 12
	println(useOutOfIf)
	useOutOfIf = 192 // want "reassigned, but never used afterwards"
	return
}
