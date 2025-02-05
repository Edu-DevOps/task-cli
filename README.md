# Example of usage

Add action:
```shell
go run main.go --add --description="prime defensivo"
Task added successfully!
```

Update action:
```shell
go run main.go --update --id=1 --description="Siuuuu" --status="in-progress"
in-progress Siuuuu
```

Delete action:
```shell
go run main.go --delete --id=1
Task 1 deleted.
```

List action:
```shell
go run main.go --list --status="To-Do"

MIRAELBICHO To-Do
Disolver las cortes To-Do
Disolver las cortes To-Do
Disolver las cortes To-Do
Disolver las cortes To-Do
Disolver las cortes To-Do
```