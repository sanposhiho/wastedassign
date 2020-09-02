# wastedassign
`wastedassign` finds wasted assignment statements

found the value ...

- reassigned, but never used afterwards
- reassigned, but reassigned soon

```
package main

func main() {
	n := 10
	println(n)

	n = 143        // want "Inefficient assignment"
	n = 13
	println(n)

	hoge := 23
	println(hoge)
	hoge = 23 // want "reassigned, but never used afterwards"
}
```
