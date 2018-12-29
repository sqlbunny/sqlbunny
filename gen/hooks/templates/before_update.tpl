	if err := {{.Var}}.doBeforeUpdateHooks(ctx); err != nil {
		return err
	}
