	if err := {{.Var}}.doBeforeInsertHooks(ctx); err != nil {
		return err
	}
