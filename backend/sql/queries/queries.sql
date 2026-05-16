-- name: CreatePlayer :one
WITH inserted_player AS (
    INSERT INTO players (username)
    VALUES ($1)
    RETURNING id
)
INSERT INTO player_stats (player_id)
SELECT id FROM inserted_player
RETURNING player_id;

-- name: GetPlayer :one
SELECT id, username FROM players
WHERE username = $1 OR id = $2
LIMIT 1;

-- name: GetPlayerStats :one
SELECT total_cash, games_played, games_won
FROM player_stats
WHERE player_id = $1
LIMIT 1;

-- name: GetPlayerActiveGame :one
SELECT games.*
FROM players
LEFT JOIN games ON games.id = players.game_id
WHERE players.id = $1
LIMIT 1;

-- name: CreateLobby :exec
WITH inserted_lobby AS (
    INSERT INTO lobbies (name, owner_id, buy_in, max_players)
    VALUES ($1, $2, $3, $4)
    RETURNING id
)
INSERT INTO lobby_players (lobby_id, player_id)
SELECT id, $2 FROM inserted_lobby
RETURNING lobby_id;

-- name: AddPlayerToLobby :exec
INSERT INTO lobby_players (lobby_id, player_id)
VALUES ($1, $2);

-- name: GetOpenLobbies :many
SELECT * FROM lobbies
WHERE status = 'WAITING' OR status = 'ACTIVE'
LIMIT $1 OFFSET $2;

