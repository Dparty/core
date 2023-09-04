# Dpary

## Configuration

Example

## Generate model from openapi

```bash
openapi-generator generate -i openapi.yaml -g go-gin-server \
  --additional-properties=packageName=model \
  --additional-properties=apiPath=model \
  -o ./controllers
```

## Generate Go SDK

```bash
openapi-generator generate -i openapi.yaml -g go \
  -o ./go-sdk
```

## Generate typescript SDK

```bash
openapi-generator generate -i openapi.yaml -g typescript-fetch -o ./boardware-cloud-ts-sdk \
   --additional-properties=npmName=boardware-cloud-ts-sdk
```

```sql
INSERT INTO accounts (id, created_at, updated_at, email, password, salt, role)
VALUES ("1681137588590612482", "2023-07-18 11:03:29.804", "2023-07-18 11:03:29.804", "chenyunda218@gmail.com", "46c24ad86e9fc5a707b495211363fa9cf0bfe7f5a0413d643e0f477ef597355ca9c4cdba87262fe3b9f6782bfa8d82e4f9c51be6843099301fdc570fe91694c6", 0x5F8038706A82641DA0BF2433079F375A, "ROOT");
```

```
GOPRIVATE=github.com/Dparty go get -u -f github.com/Dparty/core
```

## Docker run

```bash
docker run -d -it \
   --mount type=bind,source="$(pwd)"/.env.yaml,target=/.env.yaml,readonly \
   core
```

```bash
sudo docker pull chenyunda218/core
sudo docker run -dp 8080:8080 --name core --mount type=bind,source="$(pwd)"/.env.yaml,target=/app/.env.yaml chenyunda218/core
```
