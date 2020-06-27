package message

type Chain struct {
	Msg []Message
}

func GenChain(args ...Message) Chain {
	return Chain{
		Msg: args,
	}
}
