import psycopg2


CLEAR_MEDIAITEM_TABLES = ['thing_mediaitems', 'place_mediaitems', 
                'people_mediaitems', 'album_mediaitems', 
                'things', 'places', 'people', 'albums', 
                'mediaitems']
ALL_TABLES = CLEAR_MEDIAITEM_TABLES + ['users']

def after_scenario(context, scenario):
    if 'clear' in scenario.tags:
        db_conn = psycopg2.connect(
            database='carousel',
            user='carousel',
            password='carousel',
            host='localhost',
            port='5432'
        )
        for table in CLEAR_MEDIAITEM_TABLES:
            delete_table_contents(db_conn, table)
        db_conn.close()

def after_feature(context, feature):
    db_conn = psycopg2.connect(
        database='carousel',
        user='carousel',
        password='carousel',
        host='localhost',
        port='5432'
    )
    for table in ALL_TABLES:
        delete_table_contents(db_conn, table)
    db_conn.close()

def delete_table_contents(db_conn, table_name: str) -> None:
    cursor = db_conn.cursor()
    sql = f'DELETE FROM {table_name}'
    cursor.execute(sql)
    db_conn.commit()
