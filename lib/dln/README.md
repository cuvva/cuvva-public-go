# dln-go

Cuvva's library to generate, parse, and validate GB DVLA issued driving licenses.

## Usage

```go
import (
    "github.com/cuvva/dln-go"
)
```

### `Generate`

From a set of user details, generate a DLN. It cannot generate the last 3 digits though as they are not based on any user details that can be provided. Additionally, you may chose to not generate the middle name character too by passing `false` in as a 2nd argument.

```go
userDetails := dln.UserDetails{
    PersonalName: "Charles",
    FamilyName: "Bbaee",
    Sex: "M",
    BirthDate: "1975-07-01",
}

d, err := dln.Generate(userDetails, true)
if err != nil { /* do something */ }

// d = BBAEE707015C9
```

### `Validate`

Validate a DLN against a set of user details. It will return `true`/`false` depending if the details match or not. Once more, you can chose to ignore the middle name.

```go
d := "BBAEE707015C99WY"
userDetails := dln.UserDetails{
    PersonalName: "Charles",
    FamilyName: "Bbaee",
    Sex: "M",
    BirthDate: "1975-07-01",
}

isValid, err := dln.Validate(d, userDetails, true)
if err != nil { /* do something */ }

// isValid = true
```

### `Parse`

From a DLN, attempt to parse out the user details. This is not a lossless conversion, and you will loose details. Once again, you may chose to ignore any middlenames.

```go
d := "BBAEE707015C99WY"

userDetails, err := dln.Parse(d, true)
if err != nil { /* do something */ }

/*
userDetails := dln.UserDetails{
    PersonalName: "C",
    FamilyName: "Bbaee",
    Sex: "M",
    BirthDate: "1975-07-01",
}
*/
```
