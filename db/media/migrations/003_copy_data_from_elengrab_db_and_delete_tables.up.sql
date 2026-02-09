ATTACH DATABASE '${SOURCE_DB}/elengrab.db' AS source;

BEGIN;

INSERT INTO youtube_channels
SELECT *
FROM source.youtube_channels
WHERE NOT EXISTS (
    SELECT 1
    FROM youtube_channels t
    WHERE t.channel_id = source.youtube_channels.channel_id
);

INSERT INTO site_logos
SELECT *
FROM source.site_logos
WHERE NOT EXISTS (
    SELECT 1
    FROM site_logos t
    WHERE t.logo_id = source.site_logos.logo_id
);

DROP TABLE IF EXISTS source.youtube_channels;
DROP TABLE IF EXISTS source.site_logos;

COMMIT;

DETACH DATABASE source;