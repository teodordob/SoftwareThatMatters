# SoftwareThatMatters

Software that matters analysis and code

To set up go dependencies, run the following in the root directory:

```sh
go get .
```

To ingest the data, run the following in the root directory:

1. First run:
   - On Windows (requires wsl with `curl` and `sed`):
  
   ```sh
   ./init.ps1
   ```

   - On Unix systems (requires `curl` and `sed`):

   ```sh
   ./init.sh
   ```

2. Then run the following:

    ```sh
    go run .
    ```
