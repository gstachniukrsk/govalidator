package govalidator

// ExtendedWith merges two definitions, combining validators, fields, and list definitions.
func (d Definition) ExtendedWith(def Definition) Definition {
	out := Definition{}

	out.Validator = append(out.Validator, d.Validator...)
	out.Validator = append(out.Validator, def.Validator...)
	out.AcceptNotDefinedProperty = d.AcceptNotDefinedProperty || def.AcceptNotDefinedProperty
	out.AcceptExtraProperty = d.AcceptExtraProperty || def.AcceptExtraProperty

	switch {
	case d.ListOf != nil && def.ListOf != nil:
		extendedListOf := d.ListOf.ExtendedWith(*def.ListOf)
		out.ListOf = &extendedListOf
	case def.ListOf != nil:
		out.ListOf = def.ListOf
	default:
		out.ListOf = d.ListOf
	}

	f := map[string]Definition{}

	if d.Fields != nil {
		for k, v := range *d.Fields {
			f[k] = v
		}
	}

	if def.Fields != nil {
		for k, v := range *def.Fields {
			if _, ok := f[k]; ok {
				f[k] = f[k].ExtendedWith(v)
				continue
			}

			f[k] = v
		}
	}

	if len(f) > 0 {
		out.Fields = &f
	}

	return out
}
