# Lessons Learned

Random lessons learned while building this repo, so they don't have to be learned the hard way
again.


- `uci set` only works on existing sections and errors out if you provide the static fields in the `values` portion.
- UCI commands are cached client side before being committed or applied, this cache is tied to the session which
made the calls. i.e. if you do a `uci changes`, it will only show you the changes submitted on behalf of the
session you used to make the `uci changes` call.
- Manually updating pkg.go.dev:
    - curl https://sum.golang.org/lookup/github.com/daimonaslabs/go-ubus-rpc@v0.0.0-20250620052234-89352b1bdcc1
    - [How to determine pseudo-version](https://go.dev/doc/modules/version-numbers)