# miscellaneous
a repository of miscellanous things that don't need their own repos.

---

## hasher.ps1
a powershell script meant to create hashes of files in a directory. <br>
created for https://unethicalcdn.com - SHA256 hashes of files.
![hashes directory](https://github.com/unethicalteam/miscellaneous/assets/38664452/8faee606-62db-4733-8ef6-2d11be8789eb)

## url_check.go
a go script that checks an array of URLs to ensure whether they're active or dead links.  <br>
created for [lcbud](https://github.com/unethicalteam/lcbud)

***note: change the concurrent check value below to check more URLs at once in the event you have a more extensive list like 1000 URLs.***
```go
// Define Maximum Amount of Concurrent Checks
const (
	defaultMaxConcurrentChecks = 30
)
```