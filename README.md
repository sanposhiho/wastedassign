# wastedassign(English)
`wastedassign` finds wasted assignment statements

found the value ...

- reassigned, but never used afterwards
- reassigned, but reassigned soon

```
package a

func f() {
	useOutOfIf := 0 // want "wasted assignment"
	err := doHoge()
	if err != nil {
		useOutOfIf = 10 // want "wasted assignment"
		useOutOfIf = 10 // want "reassigned, but never used afterwards"

		return
	}
	
	err = doFuga() // want "reassigned, but never used afterwards"
	
	useOutOfIf = 12
	println(useOutOfIf)
	return
}
```

# wastedassign(Japanese)
`wastedassign` は無駄な代入を発見してくれる静的解析ツールです。

以下のようなstatementsを発見します。

- 代入されたがreturnまでその代入された値が使用されることはなかった
- 代入されたが代入された値が用いられることなく、別の値に変更された
