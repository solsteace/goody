-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

ALTER TABLE `produk`
    ADD COLUMN `deskripsi` TEXT NOT NULL
    AFTER `stok`;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

ALTER TABLE `produk`
    DROP COLUMN `deskripsi`;
