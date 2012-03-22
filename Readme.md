# Pest

This is a tool for testing network services.

You write a test script and run it; pest will open a
connection to the server play the test script, and
check that the returned data matches the script.

## Example

Given a file `pause.test`:

    >:put 0 0 0 1
    >:x
    <:INSERTED {.}
    s = now
    >:pause-tube default 1
    <:PAUSED
    >:reserve
    <:RESERVED {.} 1
    <:x
    ~ now - s >= 1000000000

Run it like this:

    $ pest localhost:11300 pause.test

## Documentation

Proper reference documents are forthcoming.
