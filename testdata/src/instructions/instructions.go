package instructions

func reassigned() int {
	value := 1 // want "assigned to value, but reassigned without using the value"
	value = 2
	return value
}

func readBeforeReassign() int {
	value := 1
	value += 2
	return value
}

func unusedAfterRead() {
	value := 1
	println(value)
	value = 2 // want "assigned to value, but never used afterwards"
}

func unrelatedInstructions() int {
	value := 1 // want "assigned to value, but reassigned without using the value"
	println("between assignments")
	other := 2
	println(other)
	value = 3
	return value
}

func overwriteBothBranches(condition bool) int {
	value := 1 // want "assigned to value, but reassigned without using the value"
	if condition {
		value = 2
	} else {
		value = 3
	}
	return value
}

func readOneBranch(condition bool) int {
	value := 1
	if condition {
		println(value)
	}
	value = 2
	return value
}

func usedOnBackedge(done func() bool) {
	value := 0
	for {
		println(value)
		value = 1
		if done() {
			return
		}
	}
}

func overwrittenInLoop(n int) int {
	value := 0
	for i := 0; i < n; i++ {
		value = i // want "assigned to value, but reassigned without using the value"
		value = i + 1
	}
	return value
}

func continueBackedge(n int) {
	value := 0
	for i := 0; i < n; i++ {
		println(value)
		value = i
		if i%2 == 0 {
			continue
		}
		println(value)
	}
}

func closure() func() int {
	return func() int {
		value := 1 // want "assigned to value, but reassigned without using the value"
		value = 2
		return value
	}
}

func empty() {}
