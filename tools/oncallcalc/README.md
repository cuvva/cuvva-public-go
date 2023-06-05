# oncallcalc

This CLI generates a report on on-call shifts over a specific calendar month.

It does this by splitting all days within the month to 12 hour shifts, which it then does it's best to coerce the rota into these blocks.

It additionally calculates payouts as needed.

## Requisites 
1. An accesskey from pagerduty
    1. Open [PagerDuty](https://cuvva.pagerduty.com/)
    2. Click `User Profile`
    3. Select `User Settings`
    4. Click `Create API User Token`
    5. Keep this token to be used in the next step. Optionally, add it to your 1password vault.
2. GoLang installed in your machine

## Installation

1. Set your token as an environment variable
```bash
% export PAGERDUTY_API={yourToken}
```

2. Run the installation script
```bash
% cd cmd/oncallcalc 
% go install .
```

## Commands

### `generate-report`

Generate a report on call rota payout for a month.

```bash
> go run oncallcalc.go generate-report --time "May 2023"
```
