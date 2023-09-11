package assertions

type AssertionArgs struct {
	condition bool
	message   string
}

func Assert(condition bool) {
	AssertThat(AssertionArgs{condition: condition})
}

func AssertThat(args AssertionArgs) {
	if args.message == "" {
		args.message = "assertion failed"
	}

	if !args.condition {
		panic(args.message)
	}
}
