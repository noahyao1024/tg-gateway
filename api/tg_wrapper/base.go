package tg_wrapper

type Instance struct {
	BotKey string
}

func New(botKey string) *Instance {
	i := &Instance{
		BotKey: botKey,
	}

	return i
}
