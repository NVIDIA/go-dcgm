# Contribute to the DCGM Golang Bindings

To contribute to the NVIDIA DCGM Go bindings project, sign your work. The
section below explains how.

## Build tool prerequisites

The repository uses Go 1.27.1, Task 3.53.1, Bazelisk 1.29.0, gofumpt 0.12.0,
and golangci-lint 2.13.2. Install those versions using your platform package
manager or their upstream installation instructions. Bazelisk reads the
repository's `.bazelversion` file.

You can instead open the repository in a
[Dev Containers](https://containers.dev/) compatible editor. The checked-in
development container includes DCGM headers and the pinned tools that CI uses.

Run `task versions:validate` to check repository pin consistency. Tool setup is
intentionally documented rather than performed by a privileged bootstrap
script.

## Updating DCGM Fields

When new fields are added to DCGM, you need to update the Go bindings. Follow these steps:

### 1. Update the dcgm_fields.h header file

Copy the latest `dcgm_fields.h` from the DCGM source repository:

```bash
# From the DCGM repository
cp /path/to/dcgm/dcgmlib/dcgm_fields.h pkg/dcgm/dcgm_fields.h
```

### 2. Generate Go constants

Run the code generator to update the Go field constants:

```bash
task generate
```

This will:

- Parse `pkg/dcgm/dcgm_fields.h`
- Read curated lowercase compatibility names from `pkg/dcgm/legacy_fields.csv`
- Generate `pkg/dcgm/const_fields.go` with all DCGM field constants and helper functions

### 3. Verify the generated code

Check that the generated code is correct:

```bash
task generate:check
```

This ensures the generated code is in sync with the header file.

### 4. Review the changes

Check what fields were added, removed, or modified:

```bash
git diff pkg/dcgm/const_fields.go
```

If a lowercase compatibility name needs to be added or removed, update
`pkg/dcgm/legacy_fields.csv` in the same change and regenerate the constants.

### 5. Test the changes

Run tests to ensure the bindings work correctly:

```bash
task test
task test:integration
task test:race
```

## Validate your work

All changes need to pass `task validate` and `pre-commit run --all-files`.
Changes affecting runtime behavior must also pass `task test:integration` and
`task test:race` on a qualified GPU/DCGM system.

### Setting up pre-commit

You can install pre-commit via brew, apt/dnf, or via pip:

```bash
pip install pre-commit
```

Once installed, you can run:

```bash
task pre-commit:install
pre-commit autoupdate
```

Once you complete this step, pre-commit is set up and ready to go. The pre-commit hooks
will be executed when you run `git commit`.

## Sign your work

The sign-off is a simple line at the end of the explanation for the patch. Your
signature certifies that you wrote the patch or otherwise have the right to pass
it on as an open-source patch. The rules are pretty simple: if you can certify
the below (from [developercertificate.org](http://developercertificate.org/)):

```bash
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
1 Letterman Drive
Suite D4700
San Francisco, CA, 94129

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.

Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```

Then you just add a line to every git commit message:

```bash
    Signed-off-by: Joe Smith <joe.smith@email.com>
```

Use your real name (sorry, no pseudonyms or anonymous contributions.)

If you set your `user.name` and `user.email` git configs, you can sign your
commit automatically with `git commit -s`.
