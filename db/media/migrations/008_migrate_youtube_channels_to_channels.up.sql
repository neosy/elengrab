INSERT INTO channels (
    channel_id,
    platform,
    host,
    external_id,
    channel_url,
    channel_title,
    image_url,
    image_raw,
    image_format,
    created_at,
    updated_at
)
SELECT
    printf(
        '00000000-0000-0000-0000-%012x',
        ROW_NUMBER() OVER (ORDER BY channel_id)
    ),
    'youtube',
    'www.youtube.com',
    channel_id,
    channel_url,
    channel_title,
    image_url,
    image_raw,
    image_format,
    created_at,
    updated_at
FROM youtube_channels;

DROP TABLE IF EXISTS youtube_channels;