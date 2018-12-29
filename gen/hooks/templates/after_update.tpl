	if err := {{.Var}}.doAfterUpdateHooks(ctx); err != nil {
		return err
	}
