--changeset josephspurrier:1
-- Query is missing here on purpose.
--rollback DROP TABLE user_status;

--changeset josephspurrier:2
INSERT INTO `user_status` (`id`, `status`, `created_at`, `updated_at`) VALUES
(1, 'active',   CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP),
(2, 'inactive', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP);
--rollback TRUNCATE TABLE user_status;