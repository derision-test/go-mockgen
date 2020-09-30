package command

type ConfigFunc func(*commandConfig)

func WithArgHook(argHook ArgHookFunc) ConfigFunc {
	return func(c *commandConfig) { c.argHook = argHook }
}

func WithArgValidator(argValidator ArgValidatorFunc) ConfigFunc {
	return func(c *commandConfig) { c.argValidator = argValidator }
}
