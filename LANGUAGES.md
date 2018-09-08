# License Zero CLI Language Support

`licensezero quote` and `licensezero buy` share the same subroutine for creating an inventory of packages with License Zero metadata.  At its most basic level, the subroutine recurses the current working directory, parsing and verifying every `licensezero.json` file it finds.  When it finds a `licensezero.json` file, it looks in the same directory for files that indicate a particular kind of package, like `package.json` or `setup.py`, and attempts to extract package name, package scope (user or group), and package version.

For packages installed at the system or user level, like [RubyGems](#rubygems) and [Go](#go) packages, the subroutine shells out to language-specific development tools to list dependencies, and tries to find their paths.

The relevant source files are in [`./inventory`](./inventory).

## <a id="composer">Composer</a>

- Finds dependencies by recursing the working directory.
- Reads name and version from any `composer.json` file in the same directory as any `licensezero.json`.

## <a id="go">Go</a>

- Finds dependencies by running `go list -f '{{ join .Deps "\n" }}'`.
- Finds dependency names, paths, and standard-library status by running `go list -f "$TEMPLATE" $name`.
- See <https://github.com/licensezero/cli/issues/10>

## <a id="maven">Maven</a>

- Finds dependencies by recursing the working directory.
- Reads name and version from any `pom.xml` file in the same directory as any `licensezero.json`.

## <a id="npm">npm</a>

- Finds dependencies by recursing the working directory, including `node_modules`.
- Reads name, scope, and version from any `package.json` file in the same directory as any `licensezero.json`.
- Does _not_ parse `require()` or `import` statements to find dependencies outside the working directory.

## <a id="python">Python</a>

_Incomplete Support_

- Finds dependencies by recursing the working directory.
- Reads name and version by running `python setup.py --name --version` in the same directory as any `licensezero.json`.
- See <https://github.com/licensezero/cli/issues/3>

## <a id="rubygems">RubyGems</a>

- Finds dependencies by running `bundle show`.
- Reads name and version from `bundle show` output.
- Finds dependency paths by running `bundle show --paths`.
- Does _not_ parse `require` statements to find non-Bundler dependencies.

## <a id="rust">Rust</a>

_Rudimentary Support_

- Finds dependencies by recursing the working directory.
- Does _not_ read name or version.
- Does _not_ identify packages as Rust packages.
