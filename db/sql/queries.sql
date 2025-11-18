-- name: Countries :many
SELECT id, title, code, iso_code
FROM countries
ORDER BY id;

-- name: CountryOne :one
SELECT id, title, code, iso_code
FROM countries
where id = @id;