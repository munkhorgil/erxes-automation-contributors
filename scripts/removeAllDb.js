const dbs = db.getMongo().getDBNames()

for(var i in dbs){
    const db = db.getMongo().getDB( dbs[i] );

    if (db.getName() !== "admin" && db.getName() !== "local")
    {
        print( "dropping db " + db.getName() );
        db.dropDatabase();
    }
}
