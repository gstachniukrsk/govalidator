# Migration Guide: Definition → Schema API

This guide helps you migrate from the old `Definition` API to the new `Schema` API.

## Why Migrate?

The new `Schema` API provides:
- ✅ **Cleaner syntax**: Public fields allow struct literals
- ✅ **Fluent interface**: Method chaining for better readability
- ✅ **Separated concerns**: `Required` is separate from validators
- ✅ **Convenience helpers**: `Object()` and `Array()` functions
- ✅ **Better future-proofing**: Easier to extend and maintain

## Quick Comparison

### Old Definition API
```go
validator.Definition{
    Validator: []validator.ContextValidator{
        validator.NonNullableValidator,  // ❌ Required mixed with validators
        validator.StringValidator,
        validator.MinLengthValidator(5),
    },
    Fields: &map[string]validator.Definition{ // ❌ Pointer to map
        "name": {
            Validator: []validator.ContextValidator{
                validator.NonNullableValidator,
                validator.StringValidator,
            },
        },
    },
    AcceptExtraProperty: false,      // ❌ Confusing property names
    AcceptNotDefinedProperty: true,
}
```

### New Schema API
```go
// Fluent style
NewSchema(
    IsStringValidator,
    MinLengthValidator(5),
).Required().WithFields(map[string]*Schema{  // ✅ Required separate
    "name": NewSchema(IsStringValidator).Required(),
}).WithExtra(ExtraForbid)  // ✅ Clear enum

// OR struct literal style
&Schema{
    Validators: []ContextValidator{
        IsStringValidator,
        MinLengthValidator(5),
    },
    Fields: map[string]*Schema{  // ✅ Direct map, no pointer
        "name": NewSchema(IsStringValidator).Required(),
    },
    Extra: ExtraForbid,  // ✅ Simple field
}
```

## Migration Steps

### 1. Simple Value Validators

**Before:**
```go
def := Definition{
    Validator: []ContextValidator{
        NonNullableValidator,
        StringValidator,
        MinLengthValidator(5),
    },
}
```

**After:**
```go
schema := NewSchema(
    IsStringValidator,
    MinLengthValidator(5),
).Required()

// OR
schema := &Schema{
    Validators: []ContextValidator{
        IsStringValidator,
        MinLengthValidator(5),
    },
}
schema.Required()
```

**Key changes:**
- ❌ Remove `NonNullableValidator` / `NullableValidator`
- ✅ Use `.Required()` or `.Optional()` methods
- ✅ Validator names now use `Is` prefix consistently

### 2. Object Validation

**Before:**
```go
def := Definition{
    Fields: &map[string]Definition{
        "name": {
            Validator: []ContextValidator{NonNullableValidator, StringValidator},
        },
        "age": {
            Validator: []ContextValidator{NonNullableValidator, IntValidator},
        },
    },
    AcceptExtraProperty: false,
}
```

**After (Option 1 - Fluent):**
```go
schema := Object(map[string]*Schema{
    "name": NewSchema(IsStringValidator).Required(),
    "age":  NewSchema(IsIntegerValidator).Required(),
}).WithExtra(ExtraForbid)
```

**After (Option 2 - Struct Literal):**
```go
schema := &Schema{
    Fields: map[string]*Schema{
        "name": NewSchema(IsStringValidator).Required(),
        "age":  NewSchema(IsIntegerValidator).Required(),
    },
    Extra: ExtraForbid,
}
```

**Key changes:**
- ❌ Remove pointer: `*map[string]Definition` → `map[string]*Schema`
- ✅ Use `Object()` helper for clarity
- ❌ Replace `AcceptExtraProperty: false` → ✅ `Extra: ExtraForbid`
- ❌ Replace `AcceptExtraProperty: true` → ✅ `Extra: ExtraIgnore`

### 3. Array Validation

**Before:**
```go
def := Definition{
    Validator: []ContextValidator{
        NonNullableValidator,
        IsListValidator,
    },
    ListOf: &Definition{
        Validator: []ContextValidator{
            NonNullableValidator,
            StringValidator,
        },
    },
}
```

**After (Option 1 - Array helper):**
```go
schema := Array(
    NewSchema(IsStringValidator).Required(),
).Required()
```

**After (Option 2 - Fluent):**
```go
schema := NewSchema(IsListValidator).
    WithItems(NewSchema(IsStringValidator).Required()).
    Required()
```

**After (Option 3 - Struct Literal):**
```go
schema := &Schema{
    Validators: []ContextValidator{IsListValidator},
    Items: NewSchema(IsStringValidator).Required(),
}
schema.Required()
```

**Key changes:**
- ✅ Use `Array()` helper for clarity
- ❌ Replace `ListOf` → ✅ `Items`
- ❌ Remove pointer: `*Definition` → `*Schema`

### 4. Nested Objects

**Before:**
```go
def := Definition{
    Fields: &map[string]Definition{
        "user": {
            Fields: &map[string]Definition{
                "name": {
                    Validator: []ContextValidator{
                        NonNullableValidator,
                        StringValidator,
                    },
                },
            },
        },
    },
}
```

**After:**
```go
schema := Object(map[string]*Schema{
    "user": Object(map[string]*Schema{
        "name": NewSchema(IsStringValidator).Required(),
    }),
})
```

### 5. Complex Nested Schema

**Before:**
```go
model := Definition{
    Fields: &map[string]Definition{
        "name": {
            Validator: []ContextValidator{
                NonNullableValidator,
                StringValidator,
            },
        },
        "phones": {
            Validator: []ContextValidator{
                NonNullableValidator,
                IsListValidator,
            },
            ListOf: &Definition{
                Validator: []ContextValidator{NonNullableValidator},
                Fields: &map[string]Definition{
                    "type": {
                        Validator: []ContextValidator{
                            NonNullableValidator,
                            StringValidator,
                            OneOfValidator("home", "work"),
                        },
                    },
                },
            },
        },
    },
}
```

**After:**
```go
schema := Object(map[string]*Schema{
    "name": NewSchema(IsStringValidator).Required(),
    "phones": Array(
        Object(map[string]*Schema{
            "type": NewSchema(
                IsStringValidator,
                OneOfValidator("home", "work"),
            ).Required(),
        }).Required(),
    ).Required(),
})
```

## Validator Name Changes

| Old Name | New Name |
|----------|----------|
| `StringValidator` | `IsStringValidator` |
| `IntValidator` | `IsIntegerValidator` |
| `FloatishValidator` | `FloatValidator` |
| `FloatIsLesserThanError` | `FloatTooSmallError` |
| `FloatIsGreaterThanError` | `FloatTooLargeError` |

## Field/Constant Changes

| Old | New |
|-----|-----|
| `Definition` | `Schema` |
| `ListOf` | `Items` |
| `AcceptExtraProperty: false` | `Extra: ExtraForbid` |
| `AcceptExtraProperty: true` | `Extra: ExtraIgnore` |
| `AcceptNotDefinedProperty` | *(removed - always allowed)* |
| `*map[string]Definition` | `map[string]*Schema` |

## Validation Method Changes

### Old Way
```go
v := NewBasicValidator(PathPresenter("."), SimpleErrorPresenter())
valid, errs := v.Validate(ctx, data, definition)
```

### New Way (Option 1 - Simple)
```go
valid, errs := schema.Validate(ctx, data)
```

### New Way (Option 2 - Custom Presenters)
```go
valid, errs := schema.ValidateWithPresenter(
    ctx,
    data,
    PathPresenter("."),
    DetailedErrorPresenter(),
)
```

### New Way (Option 3 - Flat List)
```go
valid, errs := schema.ValidateFlat(
    ctx,
    data,
    CombinedPresenter(".", ": "),
)
// errs is []string instead of map[string][]string
```

## Backward Compatibility

The `Definition` API is still supported! The `Schema` API converts to `Definition` internally via `ToDefinition()`:

```go
schema := NewSchema(IsStringValidator).Required()
definition := schema.ToDefinition()  // Convert for use with old code
```

## Migration Checklist

- [ ] Replace `Definition` with `Schema` or `NewSchema()`
- [ ] Remove `NonNullableValidator` and `NullableValidator` from `Validators`
- [ ] Add `.Required()` or `.Optional()` calls
- [ ] Update validator names (`StringValidator` → `IsStringValidator`, etc.)
- [ ] Replace `Fields: &map[string]Definition` with `Fields: map[string]*Schema`
- [ ] Replace `ListOf` with `Items`
- [ ] Replace `AcceptExtraProperty` with `Extra: ExtraForbid` or `Extra: ExtraIgnore`
- [ ] Consider using convenience helpers: `Object()`, `Array()`
- [ ] Update validation calls to use schema methods

## Need Help?

- Check the examples in `schema_test.go`
- See the full Schema API documentation in `schema.go`
- Original Definition API still works - migrate gradually!
