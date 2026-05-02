CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY, -- Первичный ключ
    service_name VARCHAR(255) NOT NULL, -- Название сервиса, предоставляющего подписку
    price INTEGER NOT NULL, -- Стоимость месячной подписки в рублях, копейки не учитываются
    user_id UUID NOT NULL, -- ID пользователя в формате UUID
    start_date DATE NOT NULL, -- Дата начала подписки (месяц и год)
    end_date DATE, -- Опционально дата окончания подписки
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP -- Время создания записи timestamp
)

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_service_name ON subscriptions(service_name);