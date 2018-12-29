	if err := {{.Var}}.doAfterSelectHooks(ctx); err != nil {
		return {{.Var}}, err
	}
