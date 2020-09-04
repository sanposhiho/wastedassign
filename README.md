# wastedassign
`wastedassign` finds wasted assignment statements

found the value ...

- reassigned, but never used afterwards
- reassigned, but reassigned soon

```
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
	return
}
```
