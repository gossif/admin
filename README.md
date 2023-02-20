# Admin cli for EBSI

The administration cli supports the creation of a decentralized identifier and facilitates the process of onboarding and registering the DID onto the infrastructure. The decentralized identifier uses the EBSI method to create an identifier for a legal entity. The creation of an identifier for a natural person is not supported.

## Run the cli

To see all the options: go run -tags jwx_es256k main.go

The cli verboses the http requests. To turn it off, change the WithVerbose option in the code 

## Dependecy with the ebsi package

The functions supported in this administration cli have a dependancy with the [ebsi](https://github.com/gossif/ebsi) package. 