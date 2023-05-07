import os
import psycopg2


ALL_TABLES = ['thing_mediaitems', 'place_mediaitems',
              'people_mediaitems', 'album_mediaitems',
              'things', 'places', 'people', 'albums',
              'mediaitems', 'users']

def before_feature(context, feature):
    cleanup_tables()

def after_feature(context, feature):
    cleanup_tables()

def cleanup_tables():
    # delete all rows from database
    db_conn = psycopg2.connect(
        database='smriti',
        user='smriti',
        password='smriti',
        host='localhost',
        port='5432'
    )
    for table in ALL_TABLES:
        cursor = db_conn.cursor()
        sql = f'DELETE FROM {table}'
        cursor.execute(sql)
        db_conn.commit()
    db_conn.close()
