{{- if not .NoHooks -}}
{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $modelNamePlural := .Model.Name | plural | titleCase -}}
{{- $varNamePlural := .Model.Name | plural | camelCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
func {{$varNameSingular}}BeforeInsertHook(ctx context.Context, o *{{$modelNameSingular}}) error {
	*o = {{$modelNameSingular}}{}
	return nil
}

func {{$varNameSingular}}AfterInsertHook(ctx context.Context, o *{{$modelNameSingular}}) error {
	*o = {{$modelNameSingular}}{}
	return nil
}

func {{$varNameSingular}}AfterSelectHook(ctx context.Context, o *{{$modelNameSingular}}) error {
	*o = {{$modelNameSingular}}{}
	return nil
}

func {{$varNameSingular}}BeforeUpdateHook(ctx context.Context, o *{{$modelNameSingular}}) error {
	*o = {{$modelNameSingular}}{}
	return nil
}

func {{$varNameSingular}}AfterUpdateHook(ctx context.Context, o *{{$modelNameSingular}}) error {
	*o = {{$modelNameSingular}}{}
	return nil
}

func {{$varNameSingular}}BeforeDeleteHook(ctx context.Context, o *{{$modelNameSingular}}) error {
	*o = {{$modelNameSingular}}{}
	return nil
}

func {{$varNameSingular}}AfterDeleteHook(ctx context.Context, o *{{$modelNameSingular}}) error {
	*o = {{$modelNameSingular}}{}
	return nil
}

func {{$varNameSingular}}BeforeUpsertHook(ctx context.Context, o *{{$modelNameSingular}}) error {
	*o = {{$modelNameSingular}}{}
	return nil
}

func {{$varNameSingular}}AfterUpsertHook(ctx context.Context, o *{{$modelNameSingular}}) error {
	*o = {{$modelNameSingular}}{}
	return nil
}

func test{{$modelNamePlural}}Hooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &{{$modelNameSingular}}{}
	o := &{{$modelNameSingular}}{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, {{$varNameSingular}}DBTypes, false); err != nil {
		t.Errorf("Unable to randomize {{$modelNameSingular}} object: %s", err)
	}

	Add{{$modelNameSingular}}Hook(boil.BeforeInsertHook, {{$varNameSingular}}BeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}BeforeInsertHooks = []{{$modelNameSingular}}Hook{}

	Add{{$modelNameSingular}}Hook(boil.AfterInsertHook, {{$varNameSingular}}AfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}AfterInsertHooks = []{{$modelNameSingular}}Hook{}

	Add{{$modelNameSingular}}Hook(boil.AfterSelectHook, {{$varNameSingular}}AfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}AfterSelectHooks = []{{$modelNameSingular}}Hook{}

	Add{{$modelNameSingular}}Hook(boil.BeforeUpdateHook, {{$varNameSingular}}BeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}BeforeUpdateHooks = []{{$modelNameSingular}}Hook{}

	Add{{$modelNameSingular}}Hook(boil.AfterUpdateHook, {{$varNameSingular}}AfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}AfterUpdateHooks = []{{$modelNameSingular}}Hook{}

	Add{{$modelNameSingular}}Hook(boil.BeforeDeleteHook, {{$varNameSingular}}BeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}BeforeDeleteHooks = []{{$modelNameSingular}}Hook{}

	Add{{$modelNameSingular}}Hook(boil.AfterDeleteHook, {{$varNameSingular}}AfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}AfterDeleteHooks = []{{$modelNameSingular}}Hook{}

	Add{{$modelNameSingular}}Hook(boil.BeforeUpsertHook, {{$varNameSingular}}BeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}BeforeUpsertHooks = []{{$modelNameSingular}}Hook{}

	Add{{$modelNameSingular}}Hook(boil.AfterUpsertHook, {{$varNameSingular}}AfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}AfterUpsertHooks = []{{$modelNameSingular}}Hook{}
}
{{- end}}
