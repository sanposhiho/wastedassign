package a

func f() {
	n := 10
	println(n)

	n = 143 // want "wasted assignment"
	n = 13
	println(n)

	hoge := 23
	println(hoge)
	hoge = 23 // want "reassigned, but never used afterwards"
}
