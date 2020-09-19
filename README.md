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

# Installation

```
go get -u github.com/sanposhiho/wastedassign/cmd/wastedassign
```

# Usage

```
go vet -vettool=`which wastedassign` ./...
```

# wastedassign(Japanese)
`wastedassign` は無駄な代入を発見してくれる静的解析ツールです。

以下のようなケースに役立ちます

- 無駄な代入文を省くことによる可読性アップ
- 無駄な再代入を検出することによる使用忘れの確認

また、使用しないことが明示的にわかることで、

- なぜ使用しないのか
- 関数の返り値として返す必要がそもそもないのではないか（上記画像で言うと、doFuga()はそもそもエラーを返す必要がないのではないか

などの議論が生まれるきっかけとなります。
