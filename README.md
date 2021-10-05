<p align="center">
    <a href="https://www.cuvva.com">
        <img src="https://www.cuvva.com/static/favicon.ico" width=75/>
    </a>
</p>

<h1 align="center">
    Cuvva's Go Libraries
</h1>

![](https://github.com/cuvva/cuvva-public-go/actions/workflows/go.yml/badge.svg)
![](https://github.com/cuvva/cuvva-public-go/actions/workflows/golangci-lint.yml/badge.svg)
![](https://github.com/cuvva/cuvva-public-go/actions/workflows/codeql-analysis.yml/badge.svg)

This is our public monorepo of our open-source Go libraries

# Overview

## 📚 Libraries

* `dln` - parsing UK Driver License numbers
* `ksuid` - generating sortable ids
* `slicecontains` - determining if a slice contains an object
* `vrm` - parsing UK Vehicle Registration Mark (numbers)

## 🔧 Tools

* `cdep` - our internal tool for deploying services, lambdas
* `ctxcheck` - a code analyzer, detecting misuse of ctx

## ✅ Running the tests

```bash
go test -v ./...
```

## 🤝 Contributing

We welcome contributions. Feel free to submit a PR
