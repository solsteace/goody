-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE `users` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `nama` VARCHAR(255) NOT NULL,
    `kata_sandi` VARCHAR(255) NOT NULL,
    `no_telp` VARCHAR(255) NOT NULL UNIQUE,
    `tanggal_lahir` DATE NOT NULL,
    `jenis_kelamin` VARCHAR(255) NOT NULL,
    `tentang` LONGTEXT NOT NULL,
    `pekerjaan` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL UNIQUE,
    `id_provinsi` VARCHAR(255) NOT NULL,
    `id_kota` VARCHAR(255) NOT NULL,
    `is_admin` BOOLEAN NOT NULL,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP 
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE `users`;
