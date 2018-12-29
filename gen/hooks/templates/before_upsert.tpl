	if err := {{.Var}}.doBeforeUpsertHooks(ctx); err != nil {
		return err
	}
