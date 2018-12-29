{{- $varNameSingular := .Model.Name | singular | camelCase -}}

	if len({{$varNameSingular}}AfterSelectHooks) != 0 {
		for _, obj := range {{.Var}} {
			if err := obj.doAfterSelectHooks(ctx); err != nil {
				return {{.Var}}, err
			}
		}
	}
