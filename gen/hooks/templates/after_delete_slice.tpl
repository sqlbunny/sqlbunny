{{- $varNameSingular := .Model.Name | singular | camelCase -}}

	if len({{$varNameSingular}}BeforeDeleteHooks) != 0 {
		for _, obj := range {{.Var}} {
			if err := obj.doAfterDeleteHooks(ctx); err != nil {
				return err
			}
		}
	}
