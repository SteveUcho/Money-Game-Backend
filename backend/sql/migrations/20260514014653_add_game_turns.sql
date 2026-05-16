-- +goose Up
CREATE TYPE order_side AS ENUM ('BUY', 'SELL');
CREATE TYPE option_type AS ENUM ('CALL', 'PUT');
CREATE TYPE card_effect_type AS ENUM (
  'RAISE',
  'LOWER',
  'HOLD'
);

CREATE TABLE game_turns (
    id UUID PRIMARY KEY,
    game_id UUID NOT NULL REFERENCES games(id),
    quarter_number INT NOT NULL,
    turn_number INT NOT NULL,
    started_at TIMESTAMP,
    resolved_at TIMESTAMP,
    UNIQUE(game_id, quarter_number, turn_number)
);

CREATE TABLE turn_submissions (
    id UUID PRIMARY KEY,
    turn_id UUID NOT NULL REFERENCES game_turns(id),
    game_player_id UUID NOT NULL REFERENCES player_game_states(id),
    submitted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(turn_id, game_player_id)
);

CREATE TABLE stock_orders (
    id UUID PRIMARY KEY,
    submission_id UUID NOT NULL REFERENCES turn_submissions(id),
    side order_side NOT NULL,
    quantity INT NOT NULL,
    price_at_submission INT NOT NULL
);

CREATE TABLE option_orders (
    id UUID PRIMARY KEY,
    submission_id UUID NOT NULL REFERENCES turn_submissions(id),
    contract_id UUID NOT NULL REFERENCES option_contracts(id),
    side order_side NOT NULL,
    quantity INT NOT NULL,
    premium_at_submission INT NOT NULL
);

CREATE TABLE event_card_definitions (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    effect_type card_effect_type NOT NULL,
    magnitude INT NOT NULL
);

CREATE TABLE player_event_cards (
    id UUID PRIMARY KEY,
    game_player_id UUID NOT NULL REFERENCES player_game_states(id),
    card_definition_id UUID NOT NULL REFERENCES event_card_definitions(id),
    consumed BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE played_event_cards (
    id UUID PRIMARY KEY,
    submission_id UUID NOT NULL REFERENCES turn_submissions(id),
    player_event_card_id UUID NOT NULL REFERENCES player_event_cards(id)
);

-- +goose Down
DROP TYPE order_side;
DROP TYPE option_type;
DROP TYPE card_effect_type;

DROP TABLE game_turns;
DROP TABLE turn_submissions;
DROP TABLE stock_orders;
DROP TABLE option_orders;
DROP TABLE event_card_definitions;
DROP TABLE player_event_cards;
DROP TABLE played_event_cards;
