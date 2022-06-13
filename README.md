# SoftwareThatMatters

## Instructions

This document will help one reproduce the results mentioned in *Analyzing the Criticality of NPM Packages Through a Time-Dependent Dependency Graph*

To set up go dependencies, run the following in the root directory (optional, since go should download deps automatically) :

```sh
go mod download
```

To ingest the data, run the following in the root directory:

1. First run:
   - On Windows systems (this uses `wsl`):
      - All packages: `./init.ps1`
      - If you want only n packages (example with n = 10000): `./init.ps1 10000`

   - On unix systems (requires `wget`):
      - All packages: `./init.sh`
      - If you want only n packages (example with n = 10000): `./init.sh 10000`

2. Then run the following:

    ```sh
    go run . ingest
    ```

After running these commands, you'll end up with a file containing the transformed dependencies (`out-merged.json`). This file can then be used to generate a dependency graph and interactively explore it.

To generate a graph from the acquired data and then interactively explore it (after the first two tasks), run the following:

```sh
go run . ingest
```

### License

The code's main license can be found in `LICENSE`. It also re-uses some modified gonum code, for which the license can be found in GONUM_LICENSE
