Most basic usage:

```nuget.Download("https://api.nuget.org/v3/index.json", "bootstrap", "4.0.0", "C:\\nuget_downloads")```

Very limited usage and no flexibility. Which I'm alright with.

Another usage
```nuget.GetLatestVersion("https://api.nuget.org/v3/index.json", "bootstrap")```  -- 4.0.0

Install it. 

```go get github.com/jasontconnell/nuget```
```go test -v```

   It'll create some nuget packages


