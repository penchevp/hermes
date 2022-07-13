# Hermes

Hermes is a simple RESTful service to manage customers and broadcast notifications to customers

# Environment variables

| Name | Description |
| :--- | :--- |
| DB_HOST | Database server to connect to |
| DB_PORT | Database server port |
| DB_NAME | Database to connect to |
| DB_USERNAME | User for database connection |
| DB_PASSWORD | User password for database connection |

# Dependencies
* `github.com/gorilla/mux`
* `github.com/rs/zerolog`
* `gorm.io/gorm`
* `gorm.io/driver/postgres`
* `github.com/google/uuid`

# REST API

The REST API to the Hermes is described below.

## Get customers

### Request

`GET /customers`

    curl -i --location --request GET 'http://localhost:8024/customers'

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Wed, 13 Jul 2022 17:17:29 GMT
    Content-Length: 206

    [
      {
          "ID": "00118794-6740-44eb-b635-83eda9d16b9e",
          "Name": "Jack"
      },
      {
          "ID": "23deeb17-b753-4e1e-835f-d711986e1b39",
          "Name": "Jones
      }
    ]

## Get specific customer

### Request

`GET /customers/{id}`

    curl -i --location --request GET 'http://localhost:8024/customers/00118794-6740-44eb-b635-83eda9d16b9e'

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Wed, 13 Jul 2022 17:23:45 GMT
    Content-Length: 59

    {
        "ID": "00118794-6740-44eb-b635-83eda9d16b9e",
        "Name": "Jack"
    }

## Add new customer

### Request

`POST /customers`

    curl -i --location --request POST 'http://localhost:8024/customers' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "Name": "John"
    }'

### Response

    HTTP/1.1 201 Created
    Content-Type: application/json
    Date: Wed, 13 Jul 2022 17:26:46 GMT
    Content-Length: 59

    {
        "ID": "79904621-01c1-4260-913f-83cfa9ea9ed3",
        "Name": "John"
    }

## Update existing customer

### Request

`PUT /customers/{id}`

    curl -i --location --request PUT 'http://localhost:8024/customers/67321f78-6fa7-44d6-8306-78907de0dcdc' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "Name": "Smith"
    }'

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Wed, 13 Jul 2022 17:30:17 GMT
    Content-Length: 0

## Delete specific customer

Idempotent.

### Request

`DELETE /customers/{id}`

    curl -i --location --request DELETE 'http://localhost:8024/customers/67321f78-6fa7-44d6-8306-78907de0dcdc'

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Wed, 13 Jul 2022 17:35:09 GMT
    Content-Length: 0

## Get customer notification channels

### Request

`GET /customers/{id}/notification-channels`

    curl -i --location --request GET 'http://localhost:8024/customers/00118794-6740-44eb-b635-83eda9d16b9e/notification-channels'

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Wed,
    13 Jul 2022 17: 38: 23 GMT
    Content-Length: 355

    [
        {
            "ID": 1,
            "CustomerID": "00118794-6740-44eb-b635-83eda9d16b9e",
            "Customer": {
                "ID": "00000000-0000-0000-0000-000000000000",
                "Name": ""
            },
            "NotificationChannelTypeID": "62cd7fb6-0c96-484f-9e54-2f146643a7a0",
            "NotificationChannelType": {
                "ID": "00000000-0000-0000-0000-000000000000",
                "Name": ""
            },
            "NotificationChannelLookupKey": "burned411@gmail.com",
            "ContactCustomer": false
        }
    ]

## Update customer notification channels

### Request

`POST /customers/{customer_id}/notification-channels/{notification_channel_id}`

    curl -i --location --request POST 'http://localhost:8024/customers/23deeb17-b753-4e1e-835f-d711986e1b39/notification-channels/5cbd9281-f056-48de-80f5-0e4f0d882ce8' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "contact_customer": true,
        "lookup_key": "+359888123456"
    }'

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Wed, 13 Jul 2022 17:41:27 GMT
    Content-Length: 0

## Post notification

### Request

`POST /notifications`

    curl -i --location --request POST 'http://localhost:8024/notifications' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "from": "myself",
        "text": "tarator"
    }'

### Response

    HTTP/1.1 201 Created
    Content-Type: application/json
    Date: Wed, 13 Jul 2022 17:43:39 GMT
    Content-Length: 120

    {
        "ID": "ee0fd644-de55-4fbb-aad4-54e2e0878c41",
        "CreatedAt": "2022-07-13T17:43:39.316282Z",
        "From": "myself",
        "Text": "tarator"
    }
    
## Other response codes

| Status Code | Description |
| :--- | :--- |
| 200 | `OK` |
| 201 | `CREATED` |
| 400 | `BAD REQUEST` |
| 404 | `NOT FOUND` |
| 500 | `INTERNAL SERVER ERROR` |
