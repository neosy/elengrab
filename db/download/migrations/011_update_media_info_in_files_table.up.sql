-- Migrate legacy media_info JSON:
-- 1) Rename FormatType -> formatType, Format -> format
-- 2) Move legacy video fields into videoInfo if VideoCodec != 'none'
-- 3) Remove legacy fields otherwise
UPDATE files
SET media_info =
    CASE
        WHEN json_extract(media_info, '$.VideoCodec') != 'none' THEN
            json_set(
                json_remove(
                    media_info,
                    '$.VideoCodec',
                    '$.Resolution',
                    '$.Width',
                    '$.Height',
                    '$.FormatType',
                    '$.Format'
                ),
                '$.formatType', json_extract(media_info, '$.FormatType'),
                '$.format',     json_extract(media_info, '$.Format'),
                '$.videoInfo',
                json_object(
                    'codec',      json_extract(media_info, '$.VideoCodec'),
                    'resolution', json_extract(media_info, '$.Resolution'),
                    'width',      json_extract(media_info, '$.Width'),
                    'height',     json_extract(media_info, '$.Height')
                )
            )
        ELSE
            json_set(
                json_remove(
                    media_info,
                    '$.VideoCodec',
                    '$.Resolution',
                    '$.Width',
                    '$.Height',
                    '$.FormatType',
                    '$.Format'
                ),
                '$.formatType', json_extract(media_info, '$.FormatType'),
                '$.format',     json_extract(media_info, '$.Format')
            )
    END
WHERE
    json_extract(media_info, '$.VideoCodec') IS NOT NULL
    AND json_extract(media_info, '$.videoInfo') IS NULL;