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
INSERT INTO accounts (id, created_at, updated_at, email, password,salt, role)
VALUES ("1681137588590612482", "2023-07-18 11:03:29.804", "2023-07-18 11:03:29.804", "chenyunda218@gmail.com", "d71416b14e0d3e050639e254466fe1fe7537c50e75fad21da12b8b5e1462d80488847e1a3d57d737cbf9f1046c27c09ff7ac0955c88b6ca40e5853f4c2ad0758", 0x9905071F173336CA28E579600E48B30D, "ROOT");
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
