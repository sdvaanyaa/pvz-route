# 🌐 gRPC API Contract

```proto

service OrdersService {
    // Принять заказ от курьера
    rpc AcceptOrder (AcceptOrderRequest) returns (OrderResponse);
    // Вернуть заказ курьеру
    rpc ReturnOrder (OrderIdRequest) returns (OrderResponse);
    // Выдать заказы или принять возврат клиента
    rpc ProcessOrders (ProcessOrdersRequest) returns (ProcessResult);
    // Получить список заказов
    rpc ListOrders (ListOrdersRequest) returns (OrdersList);
    // Получить список возвратов
    rpc ListReturns (ListReturnsRequest) returns (ReturnsList);
    // Получить историю изменения заказов
    rpc GetHistory (GetHistoryRequest) returns (OrderHistoryList);
    // Импорт заказов (если эта ручка делалась ранее в рамках доп заданий)
    rpc ImportOrders (ImportOrdersRequest) returns (ImportResult);
}

message AcceptOrderRequest {
    uint64 order_id = 1;
    uint64 user_id = 2;
    google.protobuf.Timestamp expires_at = 3;
    optional PackageType package = 4;
    float weight = 5;
    float price = 6;
}

message OrderIdRequest {
    uint64 order_id = 1;
}

message ProcessOrdersRequest {
    uint64 user_id = 1;
    ActionType action = 2;
    repeated uint64 order_ids = 3;
}

enum ActionType {
    // не указан
    ACTION_TYPE_UNSPECIFIED = 0;
    // выдать заказы
    ACTION_TYPE_ISSUE = 1;
    // принять возврат клиента
    ACTION_TYPE_RETURN = 2;
}

message ListOrdersRequest {
    uint64 user_id = 1;
    bool in_pvz = 2; // если true, то будут заказы для выдачи клиенту, если false, то все
    optional uint32 last_n = 3;
    optional Pagination pagination = 4;
}

message Pagination {
    uint32 page = 1;
    uint32 count_on_page = 2;
}

message ListReturnsRequest {
    Pagination pagination = 1;
}

message ImportOrdersRequest {
    repeated AcceptOrderRequest orders = 1;
}

message GetHistoryRequest {
    Pagination pagination = 1;
}

message OrderResponse {
    OrderStatus status = 1;
    uint64 order_id = 2;
}

message ProcessResult {
    repeated uint64 processed = 1;
    repeated uint64 errors = 2;
}

message OrdersList {
    repeated Order orders = 1;
    int32 total = 2;
}

message ReturnsList {
    repeated Order returns = 1;
}

message OrderHistoryList {
    repeated OrderHistory history = 1;
}

message ImportResult {
    int32 imported = 1;
    repeated uint64 errors = 2;
}

message Order {
    uint64 order_id = 1;
    uint64 user_id = 2;
    OrderStatus status = 3;
    google.protobuf.Timestamp expires_at = 4;
    float weight = 5;
    float total_price = 6;
    optional PackageType package = 7;
}

enum PackageType {
    // не указан
    PACKAGE_TYPE_UNSPECIFIED = 0;
    // пакет
    PACKAGE_TYPE_BAG = 1;
    // коробка
    PACKAGE_TYPE_BOX = 2;
    // пленка
    PACKAGE_TYPE_TAPE = 3;
    // пленка + пакет
    PACKAGE_TYPE_BAG_TAPE = 4;
    // пленка + коробка
    PACKAGE_TYPE_BOX_TAPE = 5;
}

enum OrderStatus {
    // не указан
    ORDER_STATUS_UNSPECIFIED = 0;
    // получен, ожидает выдачи клиенту
    ORDER_STATUS_EXPECTS = 1;
    // выдан клиенту
    ORDER_STATUS_ACCEPTED = 2;
    // возвращен клиентом в пвз
    ORDER_STATUS_RETURNED = 3;
    // возвращен курьеру из пвз
    ORDER_STATUS_DELETED = 4;
}

message OrderHistory {
    uint64 order_id = 1;
    OrderStatus status = 2;
    google.protobuf.Timestamp created_at = 3;
}

```

# Формат ошибок (HTTP+gRPC)

```json
{"error": { "code": "ORDER_NOT_FOUND", "message": "Order not found" }
```

*message может быть произвольным, главное использовать допустимый code*
