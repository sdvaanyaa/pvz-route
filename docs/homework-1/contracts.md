~~# 📜 CLI Contract

```
Command: help
Description: Показать список доступных команд.
Usage: help
Output:
  help
  accept-order    Принять заказ от курьера
  return-order    Вернуть заказ
  process-orders  Выдать или принять возврат
  list-orders     Получить список заказов
  list-returns    Получить список возвратов
  order-history   Получить историю заказов
  import-orders   Импорт заказов из файла

Command: accept-order
Description: Принять заказ от курьера.
Usage: accept-order --order-id <id> --user-id <id> --expires <yyyy-mm-dd>
Output (успех):
  ORDER_ACCEPTED: <order_id>
Output (ошибка):
  ERROR: <message>

Command: return-order
Description: Вернуть заказ курьеру.
Usage: return-order --order-id <id>
Output (успех):
  ORDER_RETURNED: <order_id>
Output (ошибка):
  ERROR: <message>

Command: process-orders
Description: Выдать заказы или принять возврат клиента.
Usage: process-orders --user-id <id> --action <issue|return> --order-ids <id1,id2,...>
Output (успех):
  PROCESSED: <order_id1>
  PROCESSED: <order_id2>
Output (ошибка):
  ERROR <order_id>: <reason>

Command: list-orders
Description: Получить список заказов.
Usage: list-orders --user-id <id> [--in-pvz] [--last <N>] [--page <N> --limit <M>]
Output:
  ORDER: <order_id> <user_id> <status> <expires_at>
  ...
  TOTAL: <number>

Command: list-returns
Description: Получить список возвратов.
Usage: list-returns [--page <N> --limit <M>]
Output:
  RETURN: <order_id> <user_id> <returned_at>
  ...
  PAGE: <n> LIMIT: <m>

Command: order-history
Description: Получить историю изменения заказов.
Usage: order-history
Output:
  HISTORY: <order_id> <status> <timestamp>
  ...

Command: import-orders
Description: Импорт заказов из JSON-файла.
Usage: import-orders --file <path>
Output (успех):
  IMPORTED: <count>
Output (ошибка):
  ERROR: <message>
```
# Формат ошибок
```
ERROR: <error_code>: <message>
```

# Дополнительное задание

## Добавляется новая команда:
```
Command: scroll-orders
Description: Получить список заказов по принципу бесконечной прокрутки.
Usage: scroll-orders --user-id <id> [--limit <N>]
Output:
  ORDER: <order_id> <user_id> <status> <expires_at>
  ...
  NEXT: <next_last_id>
```
### Поведение CLI:
- Запуск команды: `scroll-orders --user-id u123`
- Интерфейс работает в интерактивном цикле, ожидая команды от пользователя:
  - Ввод `next` (и нажатие Enter) — запрашивает и выводит следующую пачку заказов
  - Ввод `exit` — завершает цикл
- Значение `next_last_id` используется автоматически между итерациями
- После каждой пачки отображается строка `NEXT: <next_last_id>` — это внутреннее состояние, не требующее от пользователя указания вручную

#### Пример использования:
```
> scroll-orders --user-id u123
ORDER: ORD001 ...
ORDER: ORD002 ...
NEXT: ORD002

> next
ORDER: ORD003 ...
ORDER: ORD004 ...
NEXT: ORD004

> exit
```~~
