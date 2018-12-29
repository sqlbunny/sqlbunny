	if err := {{.Var}}.doAfterDeleteHooks(ctx); err != nil {
		return err
	}
