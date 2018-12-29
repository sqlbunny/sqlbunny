	if err := {{.Var}}.doAfterUpsertHooks(ctx); err != nil {
		return err
	}
