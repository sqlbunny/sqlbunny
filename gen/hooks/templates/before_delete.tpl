	if err := {{.Var}}.doBeforeDeleteHooks(ctx); err != nil {
		return err
	}
