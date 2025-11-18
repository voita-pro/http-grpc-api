-- name: CountryInsert :one
insert into countries (title, code, iso_code)
values (@title, @code, @iso_code)
returning id, title, code, iso_code;

-- name: CountryUpdate :one
update countries
set title = @title, code = @code, iso_code = @iso_code
where id = @id
returning id, title, code, iso_code;