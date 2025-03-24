CREATE TABLE aviation (
    id BIGSERIAL PRIMARY KEY,          -- Уникальный идентификатор самолёта (автоинкремент)
    title VARCHAR(255) NOT NULL,       -- Название авиакомпании или самолёта (строка, не может быть NULL)
    plane VARCHAR(255) NOT NULL,       -- Модель самолёта (строка, не может быть NULL)
    price NUMERIC(15, 2) NOT NULL      -- Цена (стоимость самолёта, числовое значение с двумя знаками после запятой)
);




CREATE TABLE flights (
    id BIGSERIAL PRIMARY KEY,                  -- Уникальный идентификатор рейса (автоинкремент)
    plane_id BIGINT NOT NULL,                  -- Идентификатор самолёта (внешний ключ)
    flights_number VARCHAR(255) NOT NULL,      -- Номер рейса (строка, не может быть NULL)
    departure VARCHAR(255) NOT NULL,           -- Аэропорт вылета (строка, не может быть NULL)
    arrival VARCHAR(255) NOT NULL,             -- Аэропорт прилёта (строка, не может быть NULL)
    CONSTRAINT fk_plane FOREIGN KEY (plane_id) -- Определение внешнего ключа
        REFERENCES aviation(id)               -- Ссылка на поле `id` в таблице `aviation`
        ON DELETE CASCADE                     -- При удалении записи в `aviation`, связанные записи в `flights` также удаляются
);

