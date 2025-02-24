CREATE TABLE aviation (
    id BIGSERIAL PRIMARY KEY,          -- Уникальный идентификатор самолёта (автоинкремент)
    title VARCHAR(255) NOT NULL,       -- Название авиакомпании или самолёта (строка, не может быть NULL)
    plane VARCHAR(255) NOT NULL,       -- Модель самолёта (строка, не может быть NULL)
    price NUMERIC(15, 2) NOT NULL      -- Цена (стоимость самолёта, числовое значение с двумя знаками после запятой)
);