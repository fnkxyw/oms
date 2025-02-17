# Результаты стресс-тестов для GRPC-сервиса

Я реализовал стресс-тестер, который выполняет запросы на запись с помощью метода `AcceptOrder`, а затем запросы на чтение с помощью метода `ListOrders`. В ходе тестов можно менять количество параллельных клиентов. По окончанию фиксируется время выполнения, а также создаются следующие профили:


- **`cpu.prof`** — содержит информацию о загруженности процессора во время выполнения тестов.
- **`mem.prof`** — содержит информацию о распределении памяти во время выполнения тестов.

Для **`mem.prof`** - ``curl -o mem.prof  "http://localhost:6060/debug/pprof/heap"``. Выполнить в терминале после стресс-тестов

Для **`cpu.prof`** - перед запуском стресс-тестов в терминале необоходимо выполнить:

```curl -o cpu.prof "http://localhost:6060/debug/pprof/profile?seconds=30"```
## Результаты тестирования


```bash
1000 запросов, 1 пользователь:
RPS AcceptOrder: 3907.93
Общее время выполнения на запись: 255.890125ms
RPS ListOrders: 5589.45
Общее время выполнения на чтение: 178.9085ms

Результаты тестирования:
1000 запросов 1 пользователь:

RPS AcceptOrder: 3907.93; 
Общее время выполнения на запись: 255.890125ms;
RPS ListOrders: 5589.45;
Общее время выполнения на чтение: 178.9085ms;


1000 запросов 5 пользователей:
RPS AcceptOrder: 8299.94;
Общее время выполнения на запись: 602.413833ms;
RPS ListOrders: 17787.51;
Общее время выполнения на чтение: 281.096125ms;


10000 запросов 1 пользователь:
RPS AcceptOrder: 3907.83;
Общее время выполнения на запись: 2.558966959s;
RPS ListOrders: 5983.31;
бщее время выполнения на чтение: 1.671315667s;


10000 запросов 10 пользователей:
RPS AcceptOrder: 3838.18
Общее время выполнения на запись: 26.05404075s
RPS ListOrders: 3424.92
Общее время выполнения на чтение: 29.19777525s


100000 запросов 1 пользователь:
RPS AcceptOrder: 4397.56
Общее время выполнения на запись: 22.739872583s
RPS ListOrders: 3581.56
Общее время выполнения на чтение: 27.920806875s
