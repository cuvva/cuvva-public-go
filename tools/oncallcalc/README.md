# oncallcalc

This CLI generates a report on on-call shifts over a specific calendar month.

It does this by splitting all days within the month to 12 hour shifts, which it then does it's best to coerce the rota into these blocks.

It additionally calculates payouts as needed.

## Installation

```bash
export PAGERDUTY_API=accesskey

duffleman in ~/Source/cuvva/go/cmd/oncallcalc/cli on master λ go install .
```

## Commands

### `generate-report`

Generate a report on call rota payout for a month.

```bash
> oncallcalc generate-report --time "Jan 2020"
```
