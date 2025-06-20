## Домашнее задание №2 «Выдача заказов в разной упаковке»

**Цель:** модифицировать ваш сервис, добавить возможность в ПВЗ выдавать заказы в любой из трех различных упаковок.

### Расширяем приложение по управлению ПВЗ.

Необходимо дополнить существующие команды:

1. **Принять заказ от курьера с выбором упаковки**

   Необходимо добавить на вход обязательные параметры Вес и Цена и опционально Тип упаковки (выбор из списка).
   * Товар может остаться неупакованным (например, если клиент приходит со своей коробкой)
   * Есть три вида нашей упаковки: пакет, коробка, пленка + 2 комбинации: пакет + пленка, коробка + пленка
   * Реализуйте функционал так, чтобы в будущем можно было легко добавить еще один вид упаковки
   * При выборе пакета необходимо проверять, что вес заказа меньше 10 кг, если нет, то возвращаем информативную ошибку
   * При выборе пакета стоимость заказа увеличивается на 5 рублей
   * При выборе коробки необходимо проверить, что вес заказа меньше 30 кг, если нет, то возвращаем информативную ошибку
   * При выборе коробки стоимость заказа увеличивается на 20 рублей
   * При выборе пленки дополнительных проверок не требуется
   * При выборе пленки стоимость заказа увеличивается на 1 рубль
   * При выводе старых заказов указывать вес и цену как 0

2. **Получить список заказов**

   * Необходимо расширить контакт ответа полями Тип упаковки, Вес и Цена.

### Дополнительное задание:
1. Если выполнено дополнительное задание первой недели с "бесконечным просмотром" списка заказов -
необходимо расширить контакт ответа полями Тип упаковки, Вес и Цена.

*Подробнее о новых полях смотри в файле contracts.md*

---

### Рефакторим код

Необходимо отрефакторить структуру проекта, если этого еще не было сделано в 1 домашнем задании:
1. Вынести в код в надлежащие каталоги
2. Расположить файлы по пакетам согласно выбранному подходу к группировке. Способ группировки остается на выбор автора (не монолит! :) ).
3. Добавить комментарии к экспортируемым функциям, сгенерировать документацию по проекту.
4. Для реализации бизнес-логики используйте паттерны. Использование паттерна должно упростить решение, а не привести к его дополнительной сложности!

### Создаем Makefile и настраиваем линтеры

1. Создать Makefile для проекта со списком команд:
    * **update** - Обновить зависимости проекта (go mod tidy)
    * **linter** - Запустить линтеры
    * **build** - Собрать приложение (go build)
    * **start** - Запустить приложение (запуск сбилженного бинарника)
    * **run** - Обновить зависимости, запустить линтеры, собрать и запустить приложение

2. Настроить линтеры под проект. Можно попробовать добавить следующие линтеры и поиграться настройками:
   * errcheck
   * ineffassign
   * unused
   * goconst
   * goimports
   * gocyclo
   * gocognit

### Дедлайны сдачи и проверки задания:
- 31 мая 23:59 (сдача) / 3 июня, 23:59 (проверка)
