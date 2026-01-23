-- name: GetCategory :one
SELECT * FROM categories 
WHERE category_id = $1;

-- name: GetCategoryByCode :one
SELECT * FROM categories 
WHERE category_code = $1;

-- name: ListCategories :many
SELECT * FROM categories 
ORDER BY category_id
LIMIT $1 OFFSET $2;

-- name: ListRootCategories :many
SELECT * FROM categories 
WHERE parent_category_id IS NULL
ORDER BY category_id;

-- name: ListSubCategories :many
SELECT * FROM categories 
WHERE parent_category_id = $1
ORDER BY category_id;

-- name: CreateCategory :one
INSERT INTO categories (
    category_code, name, parent_category_id, description
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: UpdateCategory :one
UPDATE categories 
SET 
    name = $2,
    parent_category_id = $3,
    description = $4
WHERE category_id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories 
WHERE category_id = $1;
