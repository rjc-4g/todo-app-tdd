# Manual

## Setup(for Windows)

```cmd
todo-app-tdd> docker compose up -d
```

## Testing

```cmd
todo-app-tdd\backend> docker compose exec app go test
```

## All Delete

```cmd
todo-app-tdd> docker compose down --rmi all --volumes --remove-orphans
```
