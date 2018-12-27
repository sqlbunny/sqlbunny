{{- if not .NoHooks -}}
{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
var {{$varNameSingular}}BeforeInsertHooks []{{$modelNameSingular}}Hook
var {{$varNameSingular}}BeforeUpdateHooks []{{$modelNameSingular}}Hook
var {{$varNameSingular}}BeforeDeleteHooks []{{$modelNameSingular}}Hook
var {{$varNameSingular}}BeforeUpsertHooks []{{$modelNameSingular}}Hook

var {{$varNameSingular}}AfterInsertHooks []{{$modelNameSingular}}Hook
var {{$varNameSingular}}AfterSelectHooks []{{$modelNameSingular}}Hook
var {{$varNameSingular}}AfterUpdateHooks []{{$modelNameSingular}}Hook
var {{$varNameSingular}}AfterDeleteHooks []{{$modelNameSingular}}Hook
var {{$varNameSingular}}AfterUpsertHooks []{{$modelNameSingular}}Hook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *{{$modelNameSingular}}) doBeforeInsertHooks(ctx context.Context) (err error) {
	for _, hook := range {{$varNameSingular}}BeforeInsertHooks {
		if err := hook(ctx, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *{{$modelNameSingular}}) doBeforeUpdateHooks(ctx context.Context) (err error) {
	for _, hook := range {{$varNameSingular}}BeforeUpdateHooks {
		if err := hook(ctx, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *{{$modelNameSingular}}) doBeforeDeleteHooks(ctx context.Context) (err error) {
	for _, hook := range {{$varNameSingular}}BeforeDeleteHooks {
		if err := hook(ctx, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *{{$modelNameSingular}}) doBeforeUpsertHooks(ctx context.Context) (err error) {
	for _, hook := range {{$varNameSingular}}BeforeUpsertHooks {
		if err := hook(ctx, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *{{$modelNameSingular}}) doAfterInsertHooks(ctx context.Context) (err error) {
	for _, hook := range {{$varNameSingular}}AfterInsertHooks {
		if err := hook(ctx, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *{{$modelNameSingular}}) doAfterSelectHooks(ctx context.Context) (err error) {
	for _, hook := range {{$varNameSingular}}AfterSelectHooks {
		if err := hook(ctx, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *{{$modelNameSingular}}) doAfterUpdateHooks(ctx context.Context) (err error) {
	for _, hook := range {{$varNameSingular}}AfterUpdateHooks {
		if err := hook(ctx, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *{{$modelNameSingular}}) doAfterDeleteHooks(ctx context.Context) (err error) {
	for _, hook := range {{$varNameSingular}}AfterDeleteHooks {
		if err := hook(ctx, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *{{$modelNameSingular}}) doAfterUpsertHooks(ctx context.Context) (err error) {
	for _, hook := range {{$varNameSingular}}AfterUpsertHooks {
		if err := hook(ctx, o); err != nil {
			return err
		}
	}

	return nil
}

// Add{{$modelNameSingular}}Hook registers your hook function for all future operations.
func Add{{$modelNameSingular}}Hook(hookPoint bunny.HookPoint, {{$varNameSingular}}Hook {{$modelNameSingular}}Hook) {
	switch hookPoint {
		case bunny.BeforeInsertHook:
			{{$varNameSingular}}BeforeInsertHooks = append({{$varNameSingular}}BeforeInsertHooks, {{$varNameSingular}}Hook)
		case bunny.BeforeUpdateHook:
			{{$varNameSingular}}BeforeUpdateHooks = append({{$varNameSingular}}BeforeUpdateHooks, {{$varNameSingular}}Hook)
		case bunny.BeforeDeleteHook:
			{{$varNameSingular}}BeforeDeleteHooks = append({{$varNameSingular}}BeforeDeleteHooks, {{$varNameSingular}}Hook)
		case bunny.BeforeUpsertHook:
			{{$varNameSingular}}BeforeUpsertHooks = append({{$varNameSingular}}BeforeUpsertHooks, {{$varNameSingular}}Hook)
		case bunny.AfterInsertHook:
			{{$varNameSingular}}AfterInsertHooks = append({{$varNameSingular}}AfterInsertHooks, {{$varNameSingular}}Hook)
		case bunny.AfterSelectHook:
			{{$varNameSingular}}AfterSelectHooks = append({{$varNameSingular}}AfterSelectHooks, {{$varNameSingular}}Hook)
		case bunny.AfterUpdateHook:
			{{$varNameSingular}}AfterUpdateHooks = append({{$varNameSingular}}AfterUpdateHooks, {{$varNameSingular}}Hook)
		case bunny.AfterDeleteHook:
			{{$varNameSingular}}AfterDeleteHooks = append({{$varNameSingular}}AfterDeleteHooks, {{$varNameSingular}}Hook)
		case bunny.AfterUpsertHook:
			{{$varNameSingular}}AfterUpsertHooks = append({{$varNameSingular}}AfterUpsertHooks, {{$varNameSingular}}Hook)
	}
}
{{- end}}
