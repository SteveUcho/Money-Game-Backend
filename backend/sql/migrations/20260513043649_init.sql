-- +goose Up
CREATE TYPE game_status AS ENUM (
  'WAITING',
  'ACTIVE',
  'FINISHED'
);

CREATE TABLE players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ory_id TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE player_stats (
    player_id UUID PRIMARY KEY REFERENCES players(id),
    total_cash BIGINT NOT NULL DEFAULT 0,
    games_played INT NOT NULL DEFAULT 0,
    games_won INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE lobbies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    ticker TEXT NOT NULL,
    owner_id UUID NOT NULL REFERENCES players(id),
    status game_status NOT NULL DEFAULT 'WAITING',
    buy_in INT NOT NULL,
    max_players INT NOT NULL DEFAULT 4,
    current_quarter INT NOT NULL DEFAULT 1,
    current_turn INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    started_at TIMESTAMP,
    ended_at TIMESTAMP
);

CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lobby_id UUID REFERENCES lobbies(id),
    name TEXT NOT NULL,
    current_price INT NOT NULL,
    current_quarter INT NOT NULL DEFAULT 1,
    current_turn INT NOT NULL DEFAULT 1
);

CREATE TABLE lobby_players (
    lobby_id UUID REFERENCES lobbies(id),
    player_id UUID REFERENCES players(id),
    PRIMARY KEY (lobby_id, player_id)
);

CREATE TABLE player_game_state (
    game_id UUID NOT NULL REFERENCES games(id),
    player_id UUID NOT NULL REFERENCES players(id),
    starting_cash INT NOT NULL,
    current_cash INT NOT NULL,
    UNIQUE(game_id, player_id),
    PRIMARY KEY (game_id, player_id)
);

-- +goose Down
DROP TABLE player_game_state;
DROP TABLE player_stats;
DROP TABLE lobby_players;
DROP TABLE players;
DROP TABLE games;
DROP TABLE lobbies;

DROP TYPE game_status;
