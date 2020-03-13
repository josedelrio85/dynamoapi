# Local DynamoDB with Docker

Playing with DynamoDB, Golang and Docker.

See this resouces:

* [Connect local dynamodb using Golang](https://gist.github.com/Tamal/02776c3e2db7eec73c001225ff52e827)
* [Testing with Dynamo Local and Go](https://medium.com/@mcleanjnathan/testing-with-dynamo-local-and-go-7b7000ef9602)
* [Configure and query DynamoDB with GoLang](https://medium.com/@spiritualcoder/step-by-step-guide-to-use-dynamodb-with-golang-cd374f159a64)
* [AWS Go code examples](https://github.com/awsdocs/aws-doc-sdk-examples/tree/master/go/example_code/dynamodb)

## Run code

```bash
go run cmd/main.go

curl -X POST -H "Content-type: application/json" http://localhost:9001/dynamo/test
```

## Docker

```bash
docker-composer up -d
```

## DynamoDB with terminal

[From this article in Medium](https://medium.com/better-programming/how-to-set-up-a-local-dynamodb-in-a-docker-container-and-perform-the-basic-putitem-getitem-38958237b968)

### create a table

```bash
aws dynamodb --endpoint-url http://localhost:8042 create-table --table-name demo-customer-info --attribute-definitions AttributeName=customerId,AttributeType=S --key-schema AttributeName=customerId,KeyType=HASH --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5
```

### create a record

```bash
aws dynamodb put-item --endpoint-url http://localhost:8042 --table-name demo-customer-info  --item '{"customerId": {"S": "1111"}, "email": {"S": "email@something.com"}}'
```

* only create a record if there’s no existing record with the same customerId

  ```bash
  --condition-expresion "attribute_not_exists(customerId)"
  ```

### retrieve a item

```bash
aws dynamodb get-item --endpoint-url http://localhost:8042 --table-name demo-customer-info --key '{"customerId": {"S":"1111"}}'
```

### update a record

```bash
aws dynamodb update-item --endpoint-url http://localhost:8042 \
  --table-name demo-customer-info \
  --key '{"customerId": {"S": "1111"}}' \
  --update-expression 'SET #email = :newEmail' \
  --expression-attribute-names '{"#email": "email"}' \
  --expression-attribute-values '{":newEmail": {"S": "newemail@somethingnew.com"}}'
```

### update a record adding more fields

```bash
aws dynamodb update-item --endpoint-url http://localhost:8042 --table-name demo-customer-info \
  --key '{"customerId": {"S": "1111"}}' \
  --update-expression 'SET #lastName = :lastName, #dateOfBirth = :dateOfBirth' \
  --expression-attribute-names '{"#lastName": "lastName", "#dateOfBirth": "dateOfBirth"}' \
  --expression-attribute-values '{":lastName": {"S": "Jones"}, ":dateOfBirth": {"S": "1985-03-01"}}' \
  --return-values ALL_NEW
  ```

### update a record with expression

```bash
aws dynamodb update-item --endpoint-url http://localhost:8042 --table-name demo-customer-info \
  --key '{"customerId": {"S": "1111"}}' \
  --update-expression 'SET #isEligibleForPromotion = :eligibility' \
  --expression-attribute-names '{"#isEligibleForPromotion": "isEligibleForPromotion", "#dateOfBirth": "dateOfBirth"}' \
  --expression-attribute-values '{":eligibility": {"BOOL": true}, ":dateFrom": {"S": "1980-01-01"}}' \
  --condition-expression "#dateOfBirth > :dateFrom" \
  --return-values ALL_NEW
```
