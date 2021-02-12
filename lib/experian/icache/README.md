# icache

Experian's "iCache" system provides access to more personal financial data than any one company should hold(!), which is very good for risk rating etc.

Authenticated using their "WASP" system - see the wasp package for that.

## Usage

Base URLs:

- Production: https://dfi.uk.experian.com/
- UAT: https://dfi.uat.uk.experian.com/

## Sources

Based on the WSDL at https://dfi.uk.experian.com/DelphiForQuotations/InteractiveWS.asmx?wsdl and several resources which Experian can provide:

- a big spreadsheet called "iCache Data Definition Document" (we have v7.4)
- a doc called "iCache Quick Start Guide" (we have v1.7)
- an example SOAP request & response (see the tests, which are roughly equivalent)

## Notes

Request marshalling currently only includes the basic structures Cuvva needs.

Response unmarshalling currently excludes:

- legacy/deprecated fields (e.g. Mosaic)
- defaulted/unavailable fields (e.g. CIFAS)
- future functionality fields (e.g. director)
- fields which repeat back inputs without anything particularly useful added
- support for generic extensions (hosted data & spare fields)
- fields which Cuvva can't currently use (e.g. vehicle, CUE, commercial)

We'd consider accepting PRs adding support for most use-cases, aside from legacy/deprecated and consistently defaulted/unavailable fields.
