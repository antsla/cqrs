Создаем заказ
`curl --request POST 'http://localhost:8080/v1/order?userId=1`

Создаем товары для заказа 
`curl --request POST http://localhost:8081/v1/goods?orderId=1`
`curl --request POST http://localhost:8081/v1/goods?orderId=1`

Получаем список:
`curl http://localhost:8082/v1/order-history?limit=2&threshold=1&offset=0`