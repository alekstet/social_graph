# social_graph [![Build Status](https://github.com/alekstet/social_graph/actions/workflows/social.yml/badge.svg)](https://github.com/alekstet/social_graph/actions/workflows/social.yml/badge.svg)


Сервис обрабатывает 2 метода:
1. PUT /social?from=1&to=2; В ответе пустое тело + статус 200 (ok). В случае пустых параметров, отрицательных, нулевых, равных - пустое тело + статус 500 (internal server error);
2. GET /social;  При наличие как минимум одной валидной коммуникации получаем:
```
{
    "matrix": [[0,1],[1,0]],
    "info": {
        "max": 1,
        "min": 1,
        "avg": 1
    }
}
```
и статус 200 (ok). В ответе содержится матрица смежности и информация по максимальному, минимальному и среднему количеству коммуникаций.

Для запуска необходимо заполнить актуальными данными файл _config.yml_


