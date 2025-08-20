-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS `addresses`;

CREATE TABLE `addresses`(
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `user_id` INT UNSIGNED NOT NULL,

    `judul_alamat` VARCHAR(255) NOT NULL,
    `nama_penerima` VARCHAR(255) NOT NULL,
    `no_telp` VARCHAR(255) NOT NULL,
    `detail_alamat` VARCHAR(255) NOT NULL,
    `updated_at` DATE NOT NULL,
    `created_at` DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (`user_id`)     
        REFERENCES `users`(`id`)
        ON DELETE CASCADE
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE `addresses`;
