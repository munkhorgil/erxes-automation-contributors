(function() {
  connection = new Mongo();

  db1 = connection.getDB('erxes');
  db2 = connection.getDB('erxes_integrations');

  db1.dropDatabase();
  db2.dropDatabase();
})();
