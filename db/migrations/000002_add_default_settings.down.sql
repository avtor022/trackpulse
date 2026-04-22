-- Remove default system settings (reverse of up migration)
DELETE FROM system_settings WHERE key IN (
    'locale',
    'web_interface_enabled',
    'ui_language',
    'db_path',
    'db_api_user',
    'db_api_password_hash',
    'hardware_com_port',
    'hardware_reader_type',
    'hardware_debounce_ms',
    'log_retention_years',
    'last_log_cleanup_date',
    'app_version',
    'schema_version'
);
