-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
ALTER TABLE `toko`
    DROP FOREIGN KEY `toko_ibfk_1`;

DROP TABLE `users`;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

CREATE TABLE `users`(
    `id` INT UNSIGNED PRIMARY KEY,
    `isAdmin` BOOLEAN
);

ALTER TABLE `toko`
    ADD CONSTRAINT `toko_ibfk_1`
    FOREIGN KEY (`id_user`) REFERENCES `users`(`id`)
