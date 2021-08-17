# LaunchDarkly Sample Go Application 

We've built a simple console application that demonstrates how LaunchDarkly's SDK works.

 Below, you'll find the basic build procedure, but for more comprehensive instructions, you can visit your [Quickstart page](https://app.launchdarkly.com/quickstart#/) or the [Go SDK reference guide](https://docs.launchdarkly.com/sdk/server-side/go).

This demo requires Go 1.14 or higher.

## Build instructions 

1. Edit `main.go` and set the value of `sdkKey` to your LaunchDarkly SDK key. If there is an existing boolean feature flag in your LaunchDarkly project that you want to evaluate, set `featureFlagKey` to the flag key. By default, it will read flags from file (local/offline mode). To read flag values from server, pass `--env` as command line arguement. e.g. `--env prod`

```go
const sdkKey = "1234567890abcdef"

const featureFlagKey = "my-flag"
```

2. On the command line, run `go build`


3. Run `./hello-go --name <UserName>`

You should see the message `"Feature flag '<flag key>' is <true/false> for <UserName> user"`.
