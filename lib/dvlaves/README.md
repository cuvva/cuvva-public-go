# dvlaves

## Testing Locally

The UAT environment can be used for testing locally, the URL is

```
https://uat.driver-vehicle-licensing.api.gov.uk/vehicle-enquiry/v1/vehicles
```

A small number of VRMs can be used in the UAT environment, each VRM below will return a corresponding status code.

```
VRM         Status Code
------------------------
AA19AAA     200
A19EEE      200
AA19PPP     200
L2WPS       200
0000        400
```
