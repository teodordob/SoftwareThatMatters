# SoftwareThatMatters

Software that matters analysis and code

To set up go dependencies, run the following in the root directory (optional, since go should download deps automatically):

```sh
go mod download
```

To ingest the data, run the following in the root directory:

1. First run:
   - On Windows (requires wsl with some unix distro installed and with `wget` and `sed`):

    ```sh
    ./init.ps1
    ```

    (If you want only n packages: [example with n = 10000])

    ```sh
    ./init.ps1 10000
    ```

   - On unix systems (requires `wget` and `sed`):

   ```sh
   ./init.sh
   ```

    (If you want only n packages: [example with n = 10000])

    ```sh
    ./init.sh 10000
    ```

2. Then run the following:

    ```sh
    go run .
    ```

After running these commands, you'll end up with a file containing the transformed dependencies (`streamedout-merged.json`). This file can then be used to generate a dependency graph.
