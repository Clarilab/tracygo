# tracygo

## How to import

> import "github.com/Clarilab/tracygo"

## How to use

Check the main.go in the example folder.
There are two middleware functions for atreugo.
UseBefore(CheckRequestID) which checks if the IDs are set if not its creates and saves them in the Uservalues.
UseAfter(WriteHeader) which writes the the userValues  into the response Header.

There is one resty middleware.
OnBeforeRequest(CheckTracingIDs) which checks if the IDs are set and writes them into the Header.
