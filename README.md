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

   - On unix systems (requires `curl` and `sed`):

   ```sh
   ./init.sh
   ```

2. Then run the following:

    ```sh
    go run .
    ```

After running these commands, you'll end up with potentially hundreds of thousands of files in the format `streamedout-i.json` and one file that merges all these in the file `streamedout-merged.json`. The latter can then be used to generate a dependency graph.
