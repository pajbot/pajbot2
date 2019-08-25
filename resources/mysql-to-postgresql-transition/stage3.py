#!/usr/bin/env python3

import pymysql.cursors
import psycopg2

mysql_connection=pymysql.connect(user='root', password='penis', database='pajbot2_test')
postgresql_connection=psycopg2.connect(user='pajbot', database='pajbot2')

for ptablename, mtablename in [('twitch_user_channel_permission', 'TwitchUserChannelPermission')]:
    with postgresql_connection.cursor() as pcursor:
        pcursor.execute('TRUNCATE TABLE '+ptablename)
        # pcursor.execute('ALTER TABLE '+ptablename + ' ALTER COLUMN permissions TYPE bit varying(64)')
        with mysql_connection.cursor() as cursor:
            cursor.execute('SELECT twitch_user_id, channel_id, CAST(permissions AS UNSIGNED) FROM '+mtablename)
            for result in cursor.fetchall():
                print('Transferring channel permission:', result)
                pcursor.execute('INSERT INTO '+ptablename+' (twitch_user_id, channel_id, permissions) VALUES(%s, %s, %s)', (result[0], result[1], result[2]))

for ptablename, mtablename in [('twitch_user_global_permission', 'TwitchUserGlobalPermission')]:
    with postgresql_connection.cursor() as pcursor:
        pcursor.execute('TRUNCATE TABLE '+ptablename)
        # pcursor.execute('ALTER TABLE '+ptablename + ' ALTER COLUMN permissions TYPE bit varying(64)')
        with mysql_connection.cursor() as cursor:
            cursor.execute('SELECT twitch_user_id, CAST(permissions AS UNSIGNED) FROM '+mtablename)
            for result in cursor.fetchall():
                print('Transferring global permission:', result)
                pcursor.execute('INSERT INTO '+ptablename+' (twitch_user_id, permissions) VALUES(%s, %s)', (result[0], result[1]))

with postgresql_connection.cursor() as pcursor:
    rename = 'ALTER TABLE moderation_action RENAME COLUMN {} TO {}'

    pcursor.execute(rename.format('channelid', 'channel_id'))
    pcursor.execute(rename.format('userid', 'user_id'))
    pcursor.execute(rename.format('targetid', 'target_id'))

with postgresql_connection.cursor() as pcursor:
    rename = 'ALTER TABLE moderation_action RENAME COLUMN {} TO {}'

    pcursor.execute('CREATE TABLE IF NOT EXISTS public.migrations (version text)')
    pcursor.execute("TRUNCATE TABLE public.migrations")
    pcursor.execute("INSERT INTO public.migrations (version) VALUES ('1')")

postgresql_connection.commit()
postgresql_connection.close()
