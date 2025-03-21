import psycopg2
import time


ALL_TABLES = ['thing_mediaitems', 'place_mediaitems',
              'people_mediaitems', 'album_mediaitems',
              'things', 'places', 'people', 'albums',
              'mediaitem_embeddings', 'mediaitem_faces', 'jobs', 'mediaitems', 'users']

def before_feature(context, feature):
    cleanup_tables()

def after_feature(context, feature):
    cleanup_tables()

def cleanup_tables():
    # delete all rows from database
    db_conn = psycopg2.connect(
        database='smriti',
        user='smritiuser',
        password='smritipass',
        host='localhost',
        port='5432'
    )
    for table in ALL_TABLES:
        cursor = db_conn.cursor()
        sql = f'DELETE FROM {table}'
        cursor.execute(sql)
        db_conn.commit()
    db_conn.close()
