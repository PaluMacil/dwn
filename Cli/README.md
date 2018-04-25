# Commandline Tool for dwn

```
cd cli
dotnet restore
dotnet user-secrets set "ConnectionStrings:Ry" "Server=127.0.0.1;Port=5432;Database=postgres;User Id=username;Password=mypw;"
```

```
dotnet publish -c Release -r win10-x64
```

```
$ENV:PATH="$ENV:PATH;$pwd\bin\Release\netcoreapp2.0\win10-x64\publish"
```