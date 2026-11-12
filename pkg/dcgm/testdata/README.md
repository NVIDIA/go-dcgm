# Public API inventory

`public_api.csv` classifies every function declared with `DCGM_PUBLIC_API` in
`dcgm_agent.h`. `TestPublicAPIInventory` checks declarations against the table
in both directions and verifies native calls for entries marked `bound`.

The inventory deliberately excludes `dcgm_test_apis.h` as a whole. DCGM labels
that header's APIs for testing, debugging, and private modules as "not officially
supported". Its use of `DCGM_PUBLIC_API` marks symbol visibility rather than a
support guarantee. Existing bindings for selected functions from that header do
not imply complete coverage of its declarations. Declarations also present in
`dcgm_agent.h` remain covered there.

A `not bound yet` entry records a public function whose Go binding is not
implemented. It requires a reason; native-call evidence is checked only for
`bound` entries.

An `excluded` function stays in the table with a reason, including deliberate
binding-scope choices such as manual policy-loop control. An exclusion does not
imply that DCGM deprecated the function unless its header says so.
