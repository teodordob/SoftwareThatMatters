# SoftwareThatMatters

Software that matters analysis and code

To set up go dependencies, run the following in the root directory (optional, since go should download deps automatically):

```sh
go mod download
```

To ingest the data, run the following in the root directory:

1. First run:
   - On Windows (requires wsl with some unix distro installed and with `curl` and `sed`):
    ```sh
    ./init.ps1
    ```
    (If you want only `n` packages:)
    ```sh
    ./init.ps1 n
    ```

   - On unix systems (requires `curl` and `sed`):

   ```sh
   ./init.sh
   ```

2. Then run the following:

    ```sh
    go run .
    ```

After running these commands, you'll end up with a file containing the transformed dependencies (`streamedout-merged.json`). This file can then be used to generate a dependency graph.
