# Manual

## Setup(for Windows)

```cmd
> docker compose up -d
```

## Testing

```cmd
> docker compose exec app go test
```

## All Delete

```cmd
> docker compose down --rmi all --volumes --remove-orphans
```
