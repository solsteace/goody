-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
ALTER TABLE `addresses` 
    RENAME COLUMN `user_id` TO `id_user`;
ALTER TABLE `addresses` RENAME `alamat`;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
ALTER TABLE `addresses`
    RENAME COLUMN `id_user` TO `user_id`;
ALTER TABLE `alamat` RENAME `addresses`;
