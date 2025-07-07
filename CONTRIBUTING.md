# Contribute to the DCGM Golang Bindings

Want to hack on the NVIDIA DCGM Golang Bindings Project? Awesome!
We only require you to sign your work, the below section describes this!

## Validate your work

All changes need to be able to pass all linting and pre-commit checks.  All tests
must pass, including `make lint-full`, `pre-commit run --all-files`, and `make test-main`

Note: There is a race in `make test-main` and it will occaisionally fail due to the race.

### Setting up pre-commit

You can install pre-commit via brew, apt/dnf, or via pip:

```bash
pip install pre-commit
```

Once installed, you can run:

```bash
make install-pre-commit
pre-commit autoupdate
```

Once you've complete this step, pre-commit is setup and ready to go.  The pre-commit hooks
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
