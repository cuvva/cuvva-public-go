# icache

Experian's "iCache" system provides access to more personal financial data than any one company should hold(!), which is very good for risk rating etc.

Authenticated using their "WASP" system - see the wasp lib for that.

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

Deprecated fields are currently excluded in the response unmarshalling.

Fields Cuvva can't currently use (e.g. CIAS, CUE, CIFAS fields etc) are not yet included in the response unmarshalling.
