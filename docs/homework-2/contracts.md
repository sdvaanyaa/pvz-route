# 📦 CLI Contract

## Изменение команды accept-order
```
Command: accept-order
Description: Принять заказ от курьера с выбором упаковки.
Usage: accept-order --order-id <id> --user-id <id> --expires <yyyy-mm-dd> [--package <bag|box|film|bag+film|box+film> --weight <float> --price <float>]
Output (успех):
  ORDER_ACCEPTED: <order_id>
  PACKAGE: <type>
  TOTAL_PRICE: <float>
Output (ошибка):
  ERROR: <message>
```

## Изменение команды scroll-orders (Если выполнено дополнительное задание первой недели)
```
Command: scroll-orders
Description: Получить список заказов по принципу бесконечной прокрутки.
Usage: scroll-orders --user-id <id> [--limit <N>]
Output:
  ORDER: <order_id> <user_id> <status> <expires_at> <package> <weight> <price>
  ...
  NEXT: <next_last_id>
```

## Изменение команды list-orders
```
Command: list-orders
Description: Получить список заказов.
Usage: list-orders --user-id <id> [--in-pvz] [--last <N>] [--page <N> --limit <M>]
Output:
  ORDER: <order_id> <user_id> <status> <expires_at> <package> <weight> <price>
  ...
  TOTAL: <number>
```

## Список допустимых вариантов упаковки
```
"package": "bag" | "box" | "film" | "bag+film" | "box+film"
```