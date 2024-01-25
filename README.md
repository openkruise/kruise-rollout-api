# kruise-rollout-api
The canonical location of the Kruise Rollout API definition and client.


## Compatibility matrix

| Kruise-Rollout-API |  Kruise-rollout  |
|--------------------|------------------|
| 0.4.1                | <= 0.4           | 
| 0.5.0                |  0.5             |

## Where does it come from?

`kruise-rollout-api` is synced from [https://github.com/openkruise/rollouts/tree/master/api](https://github.com/openkruise/rollouts/tree/master/api).
Code changes are made in that location, merged into `openkruise/rollouts` and later synced here.


### How to get it

To get the latest version, use go1.16+ and fetch using the `go get` command. For example:

```
go get github.com/openkruise/kruise-rollout-api@latest
```

To get a specific version, use go1.11+ and fetch the desired version using the `go get` command. For example:

```
go get github.com/openkruise/kruise-rollout-api@v0.4.1
```


### How to use it

please refer to the [example](examples/create-update-delete-rollout)


## Things you should NOT do

[https://github.com/openkruise/rollouts/tree/master/api](https://github.com/openkruise/rollouts/tree/master/api) is synced to here.
All changes must be made in the former. The latter is read-only.
