if err := {{.Var}}.doAfterInsertHooks(ctx); err != nil {
    return err
}
