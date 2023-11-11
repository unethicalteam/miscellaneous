# miscellaneous
a repository of miscellanous things that don't need their own repos.

---

## hasher.ps1
a powershell script meant to create hashes of files in a directory. <br>
created for https://unethicalcdn.com - SHA256 hashes of files.
![hashes directory](https://github.com/unethicalteam/miscellaneous/assets/38664452/8faee606-62db-4733-8ef6-2d11be8789eb)

## url_check.go
a go script that checks a list of URLs in `urls.json` to ensure whether they're active or dead links.  <br>
created for [lcbud](https://github.com/unethicalteam/lcbud)

***note: change these values to your preference:***
```go
    MaxRetries  = 3 // Maximum number of retries for checking a URL
    maxConcurrentChecks = 30 // Maximum URLs to check at once
```